package data_representation

import (
	"github.com/skysparq/grib2/record"
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

func (t *Template3) Parse(section record.Section5) (Definition, error) {
	return t, nil
}

func (t *Template3) GetValues(rec record.Record) ([]float32, error) {
	return nil, nil
}
