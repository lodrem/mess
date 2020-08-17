package generator

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/luncj/mess/schema"
)

type Row map[string]interface{}

type Metadata struct {
	NumRows uint  `json:"num_rows"`
	Rows    []Row `json:"rows"`

	path            string          `json:"-"`
	primaryKeys     []string        `json:"-"`
	uniqueKeys      [][]string      `json:"-"`
	primaryKeyIndex map[string]uint `json:"-"`
	uniqueKeyIndex  map[string]uint `json:"-"`
}

func Open(s *schema.Schema, path string) (*Metadata, error) {
	md := &Metadata{
		NumRows: 0,
		Rows:    []Row{},

		path: path,

		primaryKeys:     s.PrimaryKeys,
		uniqueKeys:      s.UniqueKeys,
		primaryKeyIndex: make(map[string]uint),
		uniqueKeyIndex:  make(map[string]uint),
	}

	f, err := os.Open(path)
	if err == nil {
		if err := json.NewDecoder(f).Decode(md); err != nil {
			f.Close()
			return nil, err
		}
	} else {
		f, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}

	for idx, row := range md.Rows {
		{
			pk := md.generatePK(row)
			if md.containsPK(pk) {
				return nil, fmt.Errorf("duplicated primary key %q for row index %q", pk, idx)
			}
			md.primaryKeyIndex[pk] = uint(idx)
		}

		{
			uks := md.generateUKs(row)

			for _, uk := range uks {
				if md.containsUK(uk) {
					return nil, fmt.Errorf("duplicated unique key %q for row index %q", uk, idx)
				}
				md.uniqueKeyIndex[uk] = uint(idx)
			}
		}
	}

	return md, nil
}

func (md *Metadata) containsPK(pk string) bool {
	_, found := md.primaryKeyIndex[pk]
	return found
}

func (md *Metadata) containsUKs(uks []string) bool {
	for _, uk := range uks {
		if md.containsUK(uk) {
			return true
		}
	}
	return false
}

func (md *Metadata) containsUK(uk string) bool {
	_, found := md.uniqueKeyIndex[uk]
	return found
}

func (md *Metadata) generatePK(row Row) string {
	pks := make([]string, len(md.primaryKeys))
	copy(pks, md.primaryKeys)
	for i := range pks {
		pks[i] = fmt.Sprintf("%s=%v", pks[i], row[pks[i]])
	}
	return strings.Join(pks, "&")
}

func (md *Metadata) generateUKs(row map[string]interface{}) []string {

	keys := make([]string, 0, len(md.uniqueKeys))

	for i := range md.uniqueKeys {
		uks := make([]string, len(md.uniqueKeys[i]))
		copy(uks, md.uniqueKeys[i])
		for j := range uks {
			uks[j] = fmt.Sprintf("%s=%v", uks[j], row[uks[j]])
		}
		keys = append(keys, strings.Join(uks, "&"))
	}

	return keys
}

func (md *Metadata) FindRow(pks map[string]interface{}) (uint, Row, error) {
	pk := md.generatePK(pks)
	if idx, found := md.primaryKeyIndex[pk]; !found {
		return 0, nil, fmt.Errorf("row not found")
	} else {
		return idx, md.Rows[idx], nil
	}
}

// PickRow returns random row.
func (md *Metadata) PickRow() (uint, Row, bool) {

	var indexes []uint

	for idx, row := range md.Rows {
		if row != nil {
			indexes = append(indexes, uint(idx))
		}
	}

	if len(indexes) == 0 {
		return 0, nil, false
	}

	idx := indexes[rand.Intn(len(indexes))]

	return idx, md.Rows[idx], true
}

func (md *Metadata) InsertRow(row Row) error {
	pk := md.generatePK(row)
	if md.containsPK(pk) {
		return fmt.Errorf("duplicated primary key %s for new row", pk)
	}

	uks := md.generateUKs(row)
	if md.containsUKs(uks) {
		return fmt.Errorf("duplicated unique key %q for new row", uks)
	}

	idx := md.NumRows
	md.Rows = append(md.Rows, row)
	md.NumRows++

	md.primaryKeyIndex[pk] = idx
	for _, uk := range uks {
		md.uniqueKeyIndex[uk] = idx
	}
	return nil
}

func (md *Metadata) UpdateRow(idx uint, row Row) error {
	old := md.Rows[idx]
	if old == nil {
		return fmt.Errorf("old row not found")
	}

	if err := md.DeleteRow(idx); err != nil {
		return err
	}

	md.Rows[idx] = row
	pk := md.generatePK(row)
	if md.containsPK(pk) {
		_ = md.InsertRow(old) // Rollback
		return fmt.Errorf("duplicated primary key %q for new row", pk)
	}

	uks := md.generateUKs(row)
	if md.containsUKs(uks) {
		_ = md.InsertRow(old) // Rollback
		return fmt.Errorf("duplicated unique key %q for new row", uks)
	}

	md.primaryKeyIndex[pk] = idx
	for _, uk := range uks {
		md.uniqueKeyIndex[uk] = idx
	}

	return nil
}

func (md *Metadata) DeleteRow(idx uint) error {
	row := md.Rows[idx]
	if row == nil {
		return fmt.Errorf("row not found")
	}

	pk := md.generatePK(row)
	if !md.containsPK(pk) {
		return fmt.Errorf("primary key %q not found", pk)
	}
	delete(md.primaryKeyIndex, pk)

	uks := md.generateUKs(row)
	if md.containsUKs(uks) {
		for _, uk := range uks {
			delete(md.uniqueKeyIndex, uk)
		}
	}

	md.Rows[idx] = nil

	return nil
}

func (md *Metadata) Close() error {
	var rows []Row
	for _, row := range md.Rows {
		if row != nil {
			rows = append(rows, row)
		}
	}
	md.Rows = rows
	md.NumRows = uint(len(md.Rows))

	f, err := os.Create(md.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(md)
}
