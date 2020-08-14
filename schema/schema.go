package schema

import (
	"encoding/json"
	"fmt"
	"github.com/luncj/mess/dataset"
	"math/big"
	"os"
	"sort"
)

type FieldType string

const (
	FieldTypeInt      FieldType = "int"
	FieldTypeFloat    FieldType = "float"
	FieldTypeString   FieldType = "string"
	FieldTypeJSON     FieldType = "json"
	FieldTypeDate     FieldType = "date"
	FieldTypeDateTime FieldType = "datetime"
)

type StringType string

const (
	StringTypeAscii     StringType = "ascii"
	StringTypeWord      StringType = "word"
	StringTypeSentence  StringType = "sentence"
	StringTypeParagraph StringType = "paragraph"
)

type Field struct {
	NullableRate int       `json:"nullable_rate"`
	Type         FieldType `json:"type"`

	Int struct {
		Min *big.Int `json:"min"`
		Max *big.Int `json:"max"`
	} `json:"int"`

	Float struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"float"`

	String struct {
		Type  StringType `json:"type"`
		Ascii struct {
			MinLength int `json:"min_length"`
			MaxLength int `json:"max_length"`
		} `json:"ascii"`
		Paragraph struct {
			Num int `json:"num"`
		} `json:"paragraph"`
		Sentence struct {
			Num int `json:"num"`
		} `json:"sentence"`
		Word struct {
			Num int `json:"num"`
		} `json:"word"`
	} `json:"string"`

	Date struct {
	} `json:"date"`

	Time struct {
	} `json:"time"`

	DateTime struct {
	} `json:"datetime"`
}

type Schema struct {
	Table       string           `json:"table"`
	PrimaryKeys []string         `json:"primary_keys"`
	UniqueKeys  [][]string       `json:"unique_keys"`
	Fields      map[string]Field `json:"fields"`

	keys []string
}

func FromFile(path string) (*Schema, error) {
	var s Schema
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open schema definition file: %s", err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("read schema: %s", err)
	}

	if err := s.validate(); err != nil {
		return nil, err
	}

	sort.Strings(s.PrimaryKeys)

	for i := range s.UniqueKeys {
		sort.Strings(s.UniqueKeys[i])
	}

	s.keys = KeysFromFields(s.Fields)

	return &s, nil
}

func (s *Schema) validate() error {

	if len(s.PrimaryKeys) == 0 {
		return fmt.Errorf("primary keys should not be empty")
	}

	for _, pk := range s.PrimaryKeys {
		if _, found := s.Fields[pk]; !found {
			return fmt.Errorf("primary keys %q is not defined in fields", pk)
		}
	}

	return nil
}

func (s *Schema) Keys() []string {
	return s.keys
}

func KeysFromFields(fields map[string]Field) []string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	return keys
}

func (f Field) Generate() interface{} {

	if dataset.Nullable(f.NullableRate) {
		return nil
	}

	switch f.Type {
	case FieldTypeInt:
		return dataset.IntRange(f.Int.Min, f.Int.Max)
	case FieldTypeString:
		s := f.String
		switch s.Type {
		case StringTypeAscii:
			return dataset.Ascii(s.Ascii.MinLength, s.Ascii.MaxLength)
		case StringTypeWord:
			return dataset.WordN(s.Word.Num)
		case StringTypeSentence:
			return dataset.SentenceN(s.Sentence.Num)
		case StringTypeParagraph:
			return dataset.ParagraphN(s.Paragraph.Num)
		}
	}
	panic(fmt.Sprintf("invalid field type: %s", f.Type))
}
