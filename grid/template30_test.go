package grid_test

import (
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestTemplate30(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordGridDef30)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := grid.Template30{}.Parse(rec.Grid)
	if err != nil {
		t.Fatal(err)
	}
	expected := grid.Template30{
		EarthShape:                      6,
		RadiusScaleFactor:               0,
		RadiusScaleValue:                0,
		MajorAxisScaleFactor:            0,
		MajorAxisScaleValue:             0,
		MinorAxisScaleFactor:            0,
		MinorAxisScaleValue:             0,
		PointsAlongParallel:             1799,
		PointsAlongMeridian:             1059,
		FirstLatitude:                   21138123,
		FirstLongitude:                  237280472,
		ResolutionAndComponentFlags:     8,
		LatitudeDxDySpecified:           38500000,
		LongitudeWhereLatitudeIncreases: 262500000,
		ParallelGridLength:              3000000,
		MeridianGridLength:              3000000,
		ProjectionCenterFlags:           0,
		ScanningMode:                    64,
		SecantConeFirstLatitude:         38500000,
		SecantConeSecondLatitude:        38500000,
		SouthernPoleLatitude:            0,
		SouthernPoleLongitude:           0,
	}
	if typed := template.(grid.Template30); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
}
