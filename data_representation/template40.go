package data_representation

import (
	"iter"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template40 contains the fields for Grid point data - JPEG 2000 code stream format
type Template40 struct {
	ReferenceValue         float32
	BinaryScaleFactor      int
	DecimalScaleFactor     int
	BitDepth               int
	OriginalFieldType      int
	CompressionType        int
	TargetCompressionRatio int
}

// Parse fills in the template from the provided section
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

// DecimalScale returns the decimal scale factor of the record. The decimal scale factor is used to shift the
// decimal point of a decoded value to the correct position.
func (t Template40) DecimalScale() int {
	return t.DecimalScaleFactor
}

// GetValues unpacks the record's data into the original values
func (t Template40) GetValues(_ record.Record) ([]float64, error) {
	panic("not implemented")
}

func (t Template40) ValuesIterator(_ record.Record) (iter.Seq2[int, float64], error) {
	panic(`not implemented`)
}
