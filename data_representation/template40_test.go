package data_representation_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestTemplate40(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordDataDef40)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := data_representation.Template40{}.Parse(rec.DataRepresentation)
	if err != nil {
		t.Fatal(err)
	}
	expected := data_representation.Template40{
		ReferenceValue:         1,
		BinaryScaleFactor:      0,
		DecimalScaleFactor:     2,
		BitDepth:               13,
		OriginalFieldType:      0,
		CompressionType:        0,
		TargetCompressionRatio: 255,
	}
	if typed := template.(data_representation.Template40); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
	encoded, err := json.Marshal(template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`%s`, encoded)

	err = os.WriteFile(`../.test_files/single-grib2-record-data-def-40.jp2`, rec.Data.Data, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}
