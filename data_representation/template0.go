package data_representation

import (
	"github.com/skysparq/grib2-go/record"
)

type Template0 struct {
	ReferenceValue     float32
	BinaryScaleFactor  int
	DecimalScaleFactor int
	BitsPerValue       int
	OriginalFieldType  int
}

func (t *Template0) Parse(section record.Section5) (Definition, error) {
	return t, nil
}

func (t *Template0) GetValues(rec record.Record) ([]float32, error) {
	return nil, nil
}
