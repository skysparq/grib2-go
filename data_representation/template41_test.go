package data_representation_test

import (
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestTemplate41(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordDataDef41)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := data_representation.Template41{}.Parse(rec.DataRepresentation)
	if err != nil {
		t.Fatal(err)
	}
	expected := data_representation.Template41{
		ReferenceValue:     -9990,
		BinaryScaleFactor:  0,
		DecimalScaleFactor: 1,
		BitDepth:           16,
		OriginalFieldType:  0,
	}
	if typed := template.(data_representation.Template41); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
	encoded, err := json.Marshal(template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`%s`, encoded)
}

func TestGetValuesTemplate41(t *testing.T) {
	_, r, err := test_files.Load(test_files.MrmsCompositeRefl)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	def, err := rec.DataRepresentation.Definition()
	if err != nil {
		t.Fatal(err)
	}
	values, err := def.GetValues(rec)
	if err != nil {
		t.Fatal(err)
	}

	pmin, pmax := slices.Min(values), slices.Max(values)
	pminNoSentinel := float32(999999.9)
	for _, v := range values {
		if v < pminNoSentinel && v > -99 {
			pminNoSentinel = v
		}
	}
	if expected := float32(-22); pminNoSentinel != expected {
		t.Fatalf(`expected %v but got %v`, expected, pminNoSentinel)
	}
	if expected := float32(61.5); pmax != expected {
		t.Fatalf(`expected %v but got %v`, expected, pmax)
	}
	t.Logf(`min: %v, min (no sentinels) %v, max: %v`, pmin, pminNoSentinel, pmax)
}
