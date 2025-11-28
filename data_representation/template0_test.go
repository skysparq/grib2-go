package data_representation_test

import (
	"math"
	"slices"
	"testing"

	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
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
		BitMap: record.Section6{
			Length:          6,
			BitmapIndicator: 255,
			BitmapData:      nil,
		},
		Data: record.Section7{
			Data: []byte{},
		},
	}
	result, err := template.GetValues(rec)
	if err != nil {
		t.Error(err)
	}
	if expected := []float64{1, 1, 1, 1, 1}; !slices.Equal(result, expected) {
		t.Fatalf("expected %v, got %v", expected, result)
	}
}

func TestUnpackSimpleWithBitmap(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordDataDef0Bitmap)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := data_representation.Template0{}.Parse(rec.DataRepresentation)
	if err != nil {
		t.Fatal(err)
	}
	data, err := template.GetValues(rec)
	if err != nil {
		t.Fatal(err)
	}
	if expected := 1_905_141; len(data) != expected {
		t.Fatalf(`expected %v values but got %v`, expected, len(data))
	}
	if expected := 879.9208; math.Abs(expected-data[0]) > 1e-4 {
		t.Fatalf(`expected %v but got %v`, expected, data[0])
	}
	if expected := 4174.6708; math.Abs(expected-data[243518]) > 1e-4 {
		t.Fatalf(`expected %v but got %v`, expected, data[243518])
	}
	if !math.IsNaN(data[1131644]) {
		t.Fatalf(`expected NaN but got %v`, data[1131644])
	}
	if expected := 2214.5458; math.Abs(expected-data[len(data)-1]) > 1e-4 {
		t.Fatalf(`expected %v but got %v`, expected, data[len(data)-1])
	}
}

func TestIterateConst(t *testing.T) {
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
		BitMap: record.Section6{
			Length:          6,
			BitmapIndicator: 255,
			BitmapData:      nil,
		},
		Data: record.Section7{
			Data: []byte{},
		},
	}
	expected := []float64{1, 1, 1, 1, 1}
	iterator, err := template.Values(rec)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range iterator {
		if v != expected[i] {
			t.Fatalf("expected %v at index %v, got %v", expected, i, v)
		}
	}
}
