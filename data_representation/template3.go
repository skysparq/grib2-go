package data_representation

import (
	"math"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template3 struct {
	ReferenceValue                 float32
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
	BitsUsedForGroupLengths        int
	LastGroupLength                int
	BitsUsedForScaledGroupLengths  int
	SpatialDifferenceOrder         int
	TotalSpatialDifferencingOctets int
}

func (t Template3) Parse(section record.Section5) (record.DataRepresentationDefinition, error) {
	err := checkSectionNum(section, 3)
	if err != nil {
		return t, err
	}

	data := section.DataRepresentationTemplateData
	t.ReferenceValue = u.Float32(data[0:4])
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
	t.GroupLengthReference = u.Int32(data[26:30])
	t.BitsUsedForGroupLengths = int(data[30])
	t.LastGroupLength = u.Int32(data[31:35])
	t.BitsUsedForScaledGroupLengths = int(data[35])
	t.SpatialDifferenceOrder = int(data[36])
	t.TotalSpatialDifferencingOctets = int(data[37])
	return t, nil
}

func (t Template3) GetValues(_ record.Record) ([]float32, error) {
	return nil, nil
}

func (t Template3) missingValue() float32 {
	switch t.MissingValueManagement {
	case 1:
		return float32(t.PrimaryMissingValue)
	case 2:
		return float32(t.SecondaryMissingValue)
	default:
		return float32(math.NaN())
	}
}
