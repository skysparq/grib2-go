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
	template, err := grid.Template0{}.Parse(rec.Grid)
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

func TestTemplate0Points(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordProdDef0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := grid.Template0{}.Parse(rec.Grid)
	if err != nil {
		t.Fatal(err)
	}

	points, err := template.Points()
	if err != nil {
		t.Fatal(err)
	}

	expected := template.(grid.Template0).PointsAlongMeridian * template.(grid.Template0).PointsAlongParallel
	actual := len(points.Lats)
	if actual != expected {
		t.Fatalf("expected %v latitude ponts but got %v", expected, actual)
	}

	actual = len(points.Lngs)
	if actual != expected {
		t.Fatalf("expected %v longitude points but got %v", expected, actual)
	}

	if value := points.Lats[0]; value != 90.0 {
		t.Fatalf("expected first latitude to be 90 but got %v", value)
	}

	if value := points.Lats[1038239]; value != -90.0 {
		t.Fatalf("expected first latitude to be -90 but got %v", value)
	}

	if value := points.Lngs[0]; value != -180.0 {
		t.Fatalf("expected first latitude to be -180 but got %v", value)
	}

	if value := points.Lngs[1038239]; value != 179.75 {
		t.Fatalf("expected first latitude to be 179.75 but got %v", value)
	}
}
