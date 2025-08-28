package product_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/product"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
)

func TestTemplate0(t *testing.T) {
	f, err := os.Open("../test_files/single-grib2-record.grb2")
	if err != nil {
		t.Fatal(err)
	}
	rec, err := record.ParseRecord(f, templates.Revision20120111())
	if err != nil {
		t.Fatal(err)
	}
	template, err := (&product.Template0{}).Parse(rec.ProductDefinition)
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
