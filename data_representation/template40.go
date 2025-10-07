package data_representation

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template40 struct {
	ReferenceValue         float32
	BinaryScaleFactor      int
	DecimalScaleFactor     int
	BitDepth               int
	OriginalFieldType      int
	CompressionType        int
	TargetCompressionRatio int
}

func (t Template40) Parse(section record.Section5) (record.DataRepresentationDefinition, error) {
	err := checkSectionNum(section, 40)
	if err != nil {
		return t, err
	}

	data := section.DataRepresentationTemplateData
	t.ReferenceValue = u.Float32(data[0:4])
	t.BinaryScaleFactor = u.SignAndMagnitudeInt16(data[4:6])
	t.DecimalScaleFactor = u.SignAndMagnitudeInt16(data[6:8])
	t.BitDepth = int(data[8])
	t.OriginalFieldType = int(data[9])
	t.CompressionType = int(data[10])
	t.TargetCompressionRatio = int(data[11])
	return t, nil
}

func (t Template40) GetValues(_ record.Record) ([]float32, error) {
	return nil, nil
}
