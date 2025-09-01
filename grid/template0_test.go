package grid_test

import (
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/grid"
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
	template, err := grid.Template0{}.Parse(rec.GridDefinition)
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
	if typed := template.(grid.Template0); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
}
