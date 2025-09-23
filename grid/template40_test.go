package grid_test

import (
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestTemplate40(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordGridDef40)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := grid.Template40{}.Parse(rec.GridDefinition)
	if err != nil {
		t.Fatal(err)
	}
	expected := grid.Template40{
		EarthShape:                     6,
		RadiusScaleFactor:              0,
		RadiusScaleValue:               0,
		MajorAxisScaleFactor:           0,
		MajorAxisScaleValue:            0,
		MinorAxisScaleFactor:           0,
		MinorAxisScaleValue:            0,
		PointsAlongParallel:            3072,
		PointsAlongMeridian:            1536,
		BasicAngle:                     0,
		Subdivisions:                   0,
		FirstLatitude:                  89910324,
		FirstLongitude:                 0,
		ResolutionAndComponentFlags:    48,
		LastLatitude:                   -89910324,
		LastLongitude:                  359882813,
		ParallelIncrement:              117188,
		ParallelsBetweenPoleAndEquator: 768,
		ScanningMode:                   0,
	}
	if typed := template.(grid.Template40); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
}
