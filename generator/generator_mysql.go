package generator

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/luncj/mess/dataset"

	"github.com/luncj/mess/schema"
)

type MySQLGenerator struct {
}

func NewMySQLGenerator() Generator {
	return &MySQLGenerator{}
}

func escapeKey(key string) string {
	return fmt.Sprintf("`%s`", key)
}

func escapeKeys(keys []string) []string {
	escaped := make([]string, len(keys))
	copy(escaped, keys)

	for i := range escaped {
		escaped[i] = escapeKey(escaped[i])
	}

	return escaped
}

func (m *MySQLGenerator) Generate(md *Metadata, s *schema.Schema, dml schema.DML, numRows uint) ([]string, error) {

	if numRows == 0 {
		return nil, fmt.Errorf("num of rows should be greater than zero, got %q", numRows)
	}

	var sqls = make([]string, numRows)

	var i uint = 0
	for i < numRows {
		sql, err := m.generate(md, s, dml)
		if err != nil {
			log.Fatalf("failed to generate SQL: %s", err)
		}
		sqls[i] = sql
		i++
	}

	return sqls, nil
}

func (m *MySQLGenerator) generate(md *Metadata, s *schema.Schema, dml schema.DML) (string, error) {
	switch dml {
	case schema.DMLInsert:
		return m.insertSQL(md, s)
	case schema.DMLUpdate:
		return m.updateSQL(md, s)
	case schema.DMLDelete:
		return m.deleteSQL(md, s)
	default:
		return "", fmt.Errorf("invalid DML: %q", dml)
	}
}

func (m *MySQLGenerator) normalize(t schema.FieldType, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
	case time.Time:
		switch t {
		case schema.FieldTypeDate:
			return fmt.Sprintf("'%s'", v.Format("2006-01-02"))
		case schema.FieldTypeDateTime:
			return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
		case schema.FieldTypeTime:
			return fmt.Sprintf("'%s'", v.Format("15:04:05"))
		default:
			log.Fatalf("unsupported datetime type: %s", t)
			return ""
		}
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (m *MySQLGenerator) toRow(keys []string, values []string) Row {
	n := len(keys)
	row := make(Row)
	for i := 0; i < n; i++ {
		k := keys[i]
		v := values[i]
		row[k] = v
	}
	return row
}

func (m *MySQLGenerator) insertSQL(md *Metadata, s *schema.Schema) (string, error) {
	keys := s.Keys()

	values := make([]string, len(s.Fields))
	for {
		for i, k := range keys {
			values[i] = m.normalize(s.Fields[k].Type, s.Fields[k].Generate())
		}
		row := m.toRow(keys, values)
		if err := md.InsertRow(row); err == nil {
			break
		} else {
			log.Printf("MySQL generate insert SQL: %s", err)
		}
	}

	sql := strings.Builder{}
	{
		sql.WriteString(fmt.Sprintf("INSERT INTO %s (", escapeKey(s.Table)))
		sql.WriteString(strings.Join(escapeKeys(keys), ","))
		sql.WriteString(") VALUES (")
		sql.WriteString(strings.Join(values, ","))
		sql.WriteString(");")
	}

	return sql.String(), nil
}

func (m *MySQLGenerator) updateSQL(md *Metadata, s *schema.Schema) (string, error) {

	newRow :=  make(Row)

	for len(newRow) == 0 {
		for _, k := range s.Keys() {
			if s.IsPrimaryKey(k) {
				continue
			}
			const rate = 50 // FIXME: import from schema
			if dataset.Skip(rate) {
				continue
			}
			newRow[k] = nil
		}
	}

	idx, oldRow, ok := md.PickRow()
	if !ok {
		return "", fmt.Errorf("no more data to update")
	}

	for {
		for k := range newRow {
			newRow[k] = m.normalize(s.Fields[k].Type, s.Fields[k].Generate())
		}

		for _, k := range s.PrimaryKeys {
			newRow[k] = oldRow[k]
		}

		if err := md.UpdateRow(idx, newRow); err == nil {
			break
		} else {
			log.Printf("MySQL generate update SQL: %s. It will retry.", err)
		}
	}

	sql := strings.Builder{}
	{
		sql.WriteString(fmt.Sprintf("UPDATE %s SET ", escapeKey(s.Table)))
		setFields := make([]string, 0, len(newRow)-len(s.PrimaryKeys))
		for k, v := range newRow {
			if s.IsPrimaryKey(k) {
				continue
			}
			setFields = append(setFields, fmt.Sprintf("%s = %s", escapeKey(k), v))
		}
		sql.WriteString(strings.Join(setFields, ","))
		sql.WriteString(" WHERE ")
		whereClauses := make([]string, 0, len(s.PrimaryKeys))
		for _, k := range s.PrimaryKeys {
			v := newRow[k]
			whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", escapeKey(k), v))
		}
		sql.WriteString(strings.Join(whereClauses, " AND "))
		sql.WriteString(";")
	}

	return sql.String(), nil
}

func (m *MySQLGenerator) deleteSQL(md *Metadata, s *schema.Schema) (string, error) {
	idx, row, ok := md.PickRow()
	if !ok {
		return "", fmt.Errorf("no more data to delete")
	}
	err := md.DeleteRow(idx)
	if err != nil {
		return "", err
	}

	sql := strings.Builder{}
	{
		sql.WriteString(fmt.Sprintf("DELETE FROM %s WHERE", escapeKey(s.Table)))
		whereClauses := make([]string, len(s.PrimaryKeys))
		for i, k := range s.PrimaryKeys {
			v := row[k]
			whereClauses[i] = fmt.Sprintf("%s = %s", escapeKey(k), m.normalize(s.Fields[k].Type, v))
		}
		sql.WriteString(strings.Join(whereClauses, " AND "))
		sql.WriteString(";")
	}

	return sql.String(), nil
}
