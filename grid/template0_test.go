package grid_test

import (
	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"os"
	"testing"
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
	template := &grid.Template0{}
	err = template.Parse(rec.GridDefinition)
	if err != nil {
		t.Fatal(err)
	}
	expected := grid.Template0{
		EarthShape:                  6,
		RadiusScaleFactor:           0,
		RadiusScaleValue:            0,
		MajorAxisScaleFactor:        0,
		MajorAxisScaleValue:         0,
		MinorAxisScaleFactor:        0,
		MinorAxisScaleValue:         0,
		PointsAlongParallel:         1440,
		PointsAlongMeridian:         721,
		BasicAngle:                  0,
		Subdivisions:                -1,
		FirstLatitude:               90000000,
		FirstLongitude:              0,
		ResolutionAndComponentFlags: 48,
		LastLatitude:                -90000000,
		LastLongitude:               359750000,
		ParallelIncrement:           250000,
		MeridianIncrement:           250000,
		ScanningMode:                0,
	}
	if expected != *template {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, *template)
	}
}
