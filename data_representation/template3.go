package data_representation

import (
	"fmt"
	"iter"
	"math"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template3 contains the fields for Grid Point Data - Complex Packing and Spatial Differencing
type Template3 struct {
	ReferenceValue                 float64
	BinaryScaleFactor              int
	DecimalScaleFactor             int
	BitsPerGroup                   int
	OriginalFieldType              int
	GroupSplittingMethod           int
	MissingValueManagement         int
	PrimaryMissingValue            int
	SecondaryMissingValue          int
	TotalGroups                    int
	GroupWidthReference            int
	BitsUsedForGroupWidths         int
	GroupLengthReference           int
	LengthIncrementForGroupLengths int
	LastGroupLength                int
	BitsUsedForScaledGroupLengths  int
	SpatialDifferenceOrder         int
	TotalSpatialDifferencingOctets int
}

// Parse fills in the template from the provided section
func (t Template3) Parse(section record.Section5) (record.DataRepresentationDefinition, error) {
	err := checkSectionNum(section, 3)
	if err != nil {
		return t, err
	}

	data := section.DataRepresentationTemplateData
	t.ReferenceValue = float64(u.Float32(data[0:4]))
	t.BinaryScaleFactor = u.SignAndMagnitudeInt16(data[4:6])
	t.DecimalScaleFactor = u.SignAndMagnitudeInt16(data[6:8])
	t.BitsPerGroup = int(data[8])
	t.OriginalFieldType = int(data[9])
	t.GroupSplittingMethod = int(data[10])
	t.MissingValueManagement = int(data[11])
	t.PrimaryMissingValue = u.Uint32(data[12:16])
	t.SecondaryMissingValue = u.Uint32(data[16:20])
	t.TotalGroups = u.Int32(data[20:24])
	t.GroupWidthReference = int(data[24])
	t.BitsUsedForGroupWidths = int(data[25])
	t.GroupLengthReference = u.Uint32(data[26:30])
	t.LengthIncrementForGroupLengths = int(data[30])
	t.LastGroupLength = u.Uint32(data[31:35])
	t.BitsUsedForScaledGroupLengths = int(data[35])
	t.SpatialDifferenceOrder = int(data[36])
	t.TotalSpatialDifferencingOctets = int(data[37])
	return t, nil
}

// DecimalScale returns the decimal scale factor of the record. The decimal scale factor is used to shift the
// decimal point of a decoded value to the correct position.
func (t Template3) DecimalScale() int {
	return t.DecimalScaleFactor
}

// GetValues unpacks the record's data into the original values
func (t Template3) GetValues(rec record.Record) ([]float64, error) {
	iterator, err := t.ValuesIterator(rec)
	if err != nil {
		return nil, fmt.Errorf("error getting values: %w", err)
	}
	values := make([]float64, rec.Grid.TotalPoints)
	for i, v := range iterator {
		values[i] = v
	}
	return values, nil
}

func (t Template3) ValuesIterator(rec record.Record) (iter.Seq2[int, float64], error) {
	bitmap := NewBitmapReader(rec)
	params := &ComplexParams{
		TotalPoints:              rec.Grid.TotalPoints,
		DataPoints:               rec.DataRepresentation.TotalDataPoints,
		Order:                    t.SpatialDifferenceOrder,
		SpatialOctets:            t.TotalSpatialDifferencingOctets,
		NG:                       t.TotalGroups,
		BitsPerGroup:             t.BitsPerGroup,
		BitsPerGroupWidth:        t.BitsUsedForGroupWidths,
		BitsPerScaledGroupLength: t.BitsUsedForScaledGroupLengths,
		GroupWidthReference:      t.GroupWidthReference,
		GroupLengthReference:     t.GroupLengthReference,
		GroupLengthIncrement:     t.LengthIncrementForGroupLengths,
		LastGroupLength:          t.LastGroupLength,
		Ref:                      t.ReferenceValue,
		BinaryScale:              t.BinaryScaleFactor,
		DecimalScale:             t.DecimalScaleFactor,
		MissingValueManagement:   t.MissingValueManagement,
		PrimaryMissingValue:      float64(math.Float32frombits(uint32(t.PrimaryMissingValue))),
		SecondaryMissingValue:    float64(math.Float32frombits(uint32(t.SecondaryMissingValue))),
		Bitmap:                   bitmap,
	}
	return params.UnpackComplexIterator(rec.Data.Data)
}
