package data_representation_test

import (
	"encoding/json"
	"math"
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
		LengthIncrementForGroupLengths: 1,
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
	if expected := 102744.8; data[0] != expected {
		t.Fatalf(`expected %v but got %v`, expected, data[0])
	}
	if expected := 102751.4; data[2444] != expected {
		t.Fatalf(`expected %v but got %v`, expected, data[2444])
	}
	if expected := 98455.0; data[928935] != expected {
		t.Fatalf(`expected %v but got %v`, expected, data[928935])
	}
	if expected := 102146.4; math.Abs(expected-data[len(data)-1]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[len(data)-1])
	}
}

func TestUnpackTemplate3WithBitmap(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordDataDef3Bitmap)
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
	if expected := 4.0; data[0] != expected {
		t.Fatalf(`expected %v but got %v`, expected, data[0])
	}
	if expected := 3.0; data[148235] != expected {
		t.Fatalf(`expected %v but got %v`, expected, data[148235])
	}
	if !math.IsNaN(data[277043]) {
		t.Fatalf(`expected NaN but got %v`, data[277043])
	}
	if !math.IsNaN(data[len(data)-1]) {
		t.Fatalf(`expected NaN but got %v`, data[len(data)-1])
	}
}

func TestUnpackTemplate3WithPrimaryMissingValueMngmnt(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordDataDef3PrimaryMissingValue)
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
	if expected := 0.0016; math.Abs(expected-data[0]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[0])
	}
	if expected := 0.0016; math.Abs(expected-data[1439]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[1439])
	}
	if !math.IsNaN(data[89780]) {
		t.Fatalf(`expected NaN but got %v`, data[89780])
	}
	if expected := -0.0034; math.Abs(expected-data[121019]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[121019])
	}
	if expected := -0.0044; math.Abs(expected-data[121020]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[121020])
	}
	if expected := 0.0196; math.Abs(expected-data[198924]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[198924])
	}
	if expected := 0.0006; math.Abs(expected-data[341524]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[341524])
	}
	if expected := -0.0064; math.Abs(expected-data[493748]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[493748])
	}
	if !math.IsNaN(data[493749]) {
		t.Fatalf(`expected NaN but got %v`, data[493749])
	}
	if expected := -0.0004; math.Abs(expected-data[780624]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[780624])
	}
	if expected := 0.0016; math.Abs(expected-data[len(data)-1]) > 0.00001 {
		t.Fatalf(`expected %v but got %v`, expected, data[len(data)-1])
	}
}
