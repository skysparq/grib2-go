package data_representation_test

import (
	"encoding/binary"
	"math"
	"slices"
	"testing"

	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/record"
)

func TestUnpackConst(t *testing.T) {
	template := data_representation.Template0{
		ReferenceValue:     1,
		BinaryScaleFactor:  0,
		DecimalScaleFactor: 0,
		BitsPerValue:       0,
		OriginalFieldType:  0,
	}
	rec := record.Record{
		Grid: record.Section3{
			TotalPoints: 5,
		},
		Data: record.Section7{
			Data: []byte{},
		},
	}
	result, err := template.GetValues(rec)
	if err != nil {
		t.Error(err)
	}
	if expected := []float32{1, 1, 1, 1, 1}; !slices.Equal(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestUnpackSimple(t *testing.T) {
	template := data_representation.Template0{
		ReferenceValue:     10,
		BinaryScaleFactor:  0,
		DecimalScaleFactor: 0,
		BitsPerValue:       32,
		OriginalFieldType:  0,
	}

	packed := make([]byte, 5*4)
	binary.BigEndian.PutUint32(packed[0:4], math.Float32bits(float32(1)))
	binary.BigEndian.PutUint32(packed[4:8], math.Float32bits(float32(2)))
	binary.BigEndian.PutUint32(packed[8:12], math.Float32bits(float32(3)))
	binary.BigEndian.PutUint32(packed[12:16], math.Float32bits(float32(4)))
	binary.BigEndian.PutUint32(packed[16:20], math.Float32bits(float32(5)))

	rec := record.Record{
		Grid: record.Section3{
			TotalPoints: 5,
		},
		Data: record.Section7{
			Data: packed,
		},
	}

	result, err := template.GetValues(rec)
	if err != nil {
		t.Error(err)
	}
	if expected := []float32{11, 12, 13, 14, 15}; !slices.Equal(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}
