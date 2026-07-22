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
			Templates:   singleRowGridTemplates(5),
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
	// GetValues() is normalized to row-major order with north at row 0 and west at column 0, so
	// index 0 is the north-west corner and the last index is the south-east corner. These indices
	// were chosen to match the same underlying grid points the pre-normalization test checked.
	if expected := 342.1708; math.Abs(expected-data[0]) > 1e-4 {
		t.Fatalf(`expected %v but got %v`, expected, data[0])
	}
	if expected := 4174.6708; math.Abs(expected-data[1661130]) > 1e-4 {
		t.Fatalf(`expected %v but got %v`, expected, data[1661130])
	}
	if !math.IsNaN(data[771844]) {
		t.Fatalf(`expected NaN but got %v`, data[771844])
	}
	if !math.IsNaN(data[len(data)-1]) {
		t.Fatalf(`expected NaN but got %v`, data[len(data)-1])
	}
}
