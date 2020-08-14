package generator

import "github.com/luncj/mess/schema"

type Generator interface {
	Generate(md *Metadata, s *schema.Schema, dml schema.DML, numRows uint) ([]string, error)
}
