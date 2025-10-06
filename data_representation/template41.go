package data_representation

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template41 struct {
	ReferenceValue     float32
	BinaryScaleFactor  int
	DecimalScaleFactor int
	BitDepth           int
	OriginalFieldType  int
}

func (t Template41) Parse(section record.Section5) (Definition, error) {
	err := checkSectionNum(section, 41)
	if err != nil {
		return t, err
	}

	data := section.DataRepresentationTemplateData
	t.ReferenceValue = u.Float32(data[0:4])
	t.BinaryScaleFactor = u.SignAndMagnitudeInt16(data[4:6])
	t.DecimalScaleFactor = u.SignAndMagnitudeInt16(data[6:8])
	t.BitDepth = int(data[8])
	t.OriginalFieldType = int(data[9])
	return t, nil
}

func (t Template41) GetValues(rec record.Record) ([]float32, error) {
	return nil, nil
}
