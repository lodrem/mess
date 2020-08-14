package schema

import "fmt"

type DML string

const (
	DMLInsert DML = "insert"
	DMLUpdate DML = "update"
	DMLDelete DML = "delete"
)

func (d DML) String() string {
	return string(d)
}

func DMLFromString(s string) (DML, error) {
	switch s {
	case DMLInsert.String():
		return DMLInsert, nil
	case DMLUpdate.String():
		return DMLUpdate, nil
	case DMLDelete.String():
		return DMLDelete, nil
	default:
		return "", fmt.Errorf("invalid dml %q, required (%q / %q / %q)", s, DMLInsert, DMLUpdate, DMLDelete)
	}
}
