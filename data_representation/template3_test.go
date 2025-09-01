package data_representation_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestTemplate3(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordProdDef0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := data_representation.Template3{}.Parse(rec.DataRepresentation)
	if err != nil {
		t.Fatal(err)
	}
	expected := data_representation.Template3{
		ReferenceValue:                 952124,
		BinaryScaleFactor:              1,
		DecimalScaleFactor:             1,
		BitsPerGroup:                   15,
		OriginalFieldType:              0,
		GroupSplittingMethod:           1,
		MissingValueManagement:         0,
		PrimaryMissingValue:            1649987994,
		SecondaryMissingValue:          4294967295,
		TotalGroups:                    28896,
		GroupWidthReference:            0,
		BitsUsedForGroupWidths:         4,
		GroupLengthReference:           1,
		BitsUsedForGroupLengths:        1,
		LastGroupLength:                41,
		BitsUsedForScaledGroupLengths:  7,
		SpatialDifferenceOrder:         2,
		TotalSpatialDifferencingOctets: 3,
	}
	if typed := template.(data_representation.Template3); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
	encoded, err := json.Marshal(template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`%s`, encoded)
}

func TestUnpackTemplate3(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordProdDef0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := data_representation.Template3{}.Parse(rec.DataRepresentation)
	if err != nil {
		t.Fatal(err)
	}
	data, err := template.GetValues(rec)
	if err != nil {
		t.Fatal(err)
	}
	if expected := 1_038_240; len(data) != expected {
		t.Fatalf(`expected %v values but got %v`, expected, len(data))
	}
	if expected := float32(102744.8); data[0] != expected {
		t.Fatalf(`expected %v but got %v`, expected, data[0])
	}
	if expected := float32(102751.4); data[2444] != expected {
		t.Fatalf(`expected %v but got %v`, expected, data[2444])
	}
	if expected := float32(98455.0); data[928935] != expected {
		t.Fatalf(`expected %v but got %v`, expected, data[928935])
	}
}
