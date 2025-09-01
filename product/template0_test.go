package product_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/product"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestTemplate0(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordProdDef0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := product.Template0{}.Parse(rec.ProductDefinition)
	if err != nil {
		t.Fatal(err)
	}
	expected := product.Template0{
		DefinitionHeader: product.DefinitionHeader{
			ParameterCategory: 3,
			ParameterNumber:   1,
		},
		GeneratingProcessType:       2,
		BackgroundIdentifier:        0,
		GeneratingProcessIdentifier: 96,
		HoursAfterReference:         0,
		MinutesAfterReference:       0,
		UnitOfTimeRange:             1,
		ForecastTimeInUnits:         22,
		FirstSurfaceType:            101,
		FirstSurfaceScaleFactor:     0,
		FirstSurfaceScaleValue:      0,
		SecondSurfaceType:           255,
		SecondSurfaceScaleFactor:    0,
		SecondSurfaceScaleValue:     0,
	}
	if typed := template.(product.Template0); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
	encoded, err := json.Marshal(template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`%s`, encoded)
}
