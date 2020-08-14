package generator

import (
	"fmt"
	"log"
	"strings"

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

func (m *MySQLGenerator) normalize(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("'%s'", v)
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
			values[i] = m.normalize(s.Fields[k].Generate())
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
	pks := make(map[string]interface{})
	for _, pk := range s.PrimaryKeys {
		pks[pk] = true
	}

	keys := make([]string, len(s.PrimaryKeys))
	copy(keys, s.PrimaryKeys)
	for len(keys) == len(pks) {
		for _, k := range s.Keys() {
			if _, found := pks[k]; found {
				continue
			}
			const rate = 50 // FIXME: import from schema
			if dataset.Skip(rate) {
				continue
			}
			keys = append(keys, k)
		}
	}

	idx, oldRow, ok := md.PickRow()
	if !ok {
		return "", fmt.Errorf("no more data to update")
	} else {
		for pk := range pks {
			pks[pk] = oldRow[pk]
		}
	}
	values := make([]string, len(s.Fields))
	for {
		for i, k := range keys {
			values[i] = m.normalize(s.Fields[k].Generate())
		}
		row := m.toRow(keys, values)
		if err := md.UpdateRow(idx, row); err == nil {
			break
		} else {
			log.Printf("MySQL generate update SQL: %s", err)
		}
	}

	sql := strings.Builder{}
	{
		sql.WriteString(fmt.Sprintf("UPDATE %s SET ", escapeKey(s.Table)))
		setFields := make([]string, 0, len(keys)-len(pks))
		for i := range keys {
			if _, found := pks[keys[i]]; found {
				continue
			}
			setFields = append(setFields, fmt.Sprintf("%s = %s", escapeKey(keys[i]), values[i]))
		}
		sql.WriteString(strings.Join(setFields, ","))
		sql.WriteString(" WHERE ")
		whereClauses := make([]string, 0, len(pks))
		for k, v := range pks {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = %s", escapeKey(k), v))
		}
		sql.WriteString(strings.Join(whereClauses, ", "))
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
			whereClauses[i] = fmt.Sprintf("%s = %s", escapeKey(k), m.normalize(v))
		}
		sql.WriteString(strings.Join(whereClauses, ", "))
		sql.WriteString(";")
	}

	return sql.String(), nil
}
