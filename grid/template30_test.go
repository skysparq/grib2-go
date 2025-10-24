package grid_test

import (
	"math"
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
		EarthShape:                  6,
		RadiusScaleFactor:           0,
		RadiusScaleValue:            0,
		MajorAxisScaleFactor:        0,
		MajorAxisScaleValue:         0,
		MinorAxisScaleFactor:        0,
		MinorAxisScaleValue:         0,
		Nx:                          1799,
		Ny:                          1059,
		La1:                         21138123,
		Lo1:                         237280472,
		ResolutionAndComponentFlags: 8,
		LaD:                         38500000,
		LoV:                         262500000,
		Dx:                          3000000,
		Dy:                          3000000,
		ProjectionCenterFlags:       0,
		ScanningMode:                64,
		Latin1:                      38500000,
		Latin2:                      38500000,
		SouthernPoleLatitude:        0,
		SouthernPoleLongitude:       0,
	}
	if typed := template.(grid.Template30); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
}

func TestTemplate30Points(t *testing.T) {
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

	points, err := template.Points()
	if err != nil {
		t.Fatal(err)
	}

	if expected := -122.71953; math.Abs(expected-points.Lngs[0]) > 0.001 {
		t.Fatalf("expected first longitude to be %v but got %v", expected, points.Lngs[0])
	}
	if expected := 21.138123; math.Abs(expected-points.Lats[0]) > 0.001 {
		t.Fatalf("expected first latitude to be %v but got %v", expected, points.Lats[0])
	}

	t.Logf("lngs: %+v\nlats: %+v", points.Lngs[:100], points.Lats[:100])
}
