package product_test

import (
	"github.com/skysparq/grib2-go/product"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"os"
	"testing"
)

func TestTemplate8(t *testing.T) {
	f, err := os.Open("../test_files/single-grib2-record-prod-def-8.grb2")
	if err != nil {
		t.Fatal(err)
	}
	rec, err := record.ParseRecord(f, templates.Revision20120111())
	if err != nil {
		t.Fatal(err)
	}
	template := &product.Template0{}
	err = template.Parse(rec.ProductDefinition)
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
	if expected != *template {
		t.Fatalf(`expected %+v but got %+v`, expected, template)
	}
}
