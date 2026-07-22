package grid_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/projections"
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

	// This grid's raw scan order starts at the south-west corner (La1/Lo1), but Points() normalizes
	// to row-major order with north at row 0 and west at column 0, so index 0 is the north-west corner.
	if expected := -134.095480; math.Abs(expected-points.Lngs[0]) > 0.001 {
		t.Fatalf("expected first longitude to be %v but got %v", expected, points.Lngs[0])
	}
	if expected := 47.838623; math.Abs(expected-points.Lats[0]) > 0.001 {
		t.Fatalf("expected first latitude to be %v but got %v", expected, points.Lats[0])
	}

	t.Logf("lngs: %+v\nlats: %+v", points.Lngs[:100], points.Lats[:100])
}

// TestTemplate30SrsWktMatchesPoints verifies that SrsWkt() actually describes the same projection
// Points() used to generate the grid: it builds a Lambert Conformal Conic projection purely from
// parameters parsed out of the WKT string, forward-projects the (already normalized) generated
// lat/lng points through it, and checks the result is an exactly regular Dx x Dy grid in the
// projected plane - which is only possible if the WKT's parameters match the grid's.
func TestTemplate30SrsWktMatchesPoints(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordGridDef30)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	definition, err := grid.Template30{}.Parse(rec.Grid)
	if err != nil {
		t.Fatal(err)
	}
	template := definition.(grid.Template30)

	points, err := template.Points()
	if err != nil {
		t.Fatal(err)
	}
	wkt, err := template.SrsWkt()
	if err != nil {
		t.Fatal(err)
	}

	proj := projections.NewLambertConformalConic(
		wktSpheroidRadius(t, wkt),
		0, // the WKT describes a sphere, not an ellipsoid
		wktParam(t, wkt, "Latitude_Of_Origin"),
		wktParam(t, wkt, "Central_Meridian"),
		wktParam(t, wkt, "Standard_Parallel_1"),
		wktParam(t, wkt, "Standard_Parallel_2"),
		wktParam(t, wkt, "False_Easting"),
		wktParam(t, wkt, "False_Northing"),
	)

	// Points() is normalized to row-major order, so index 0 is the north-west corner, index 1 is
	// one column east, and index Ni is one row south.
	ni := template.XVals()
	dxMeters := float64(template.Dx) * 1e-3
	dyMeters := float64(template.Dy) * 1e-3

	x00, y00 := proj.Forward(points.Lats[0], points.Lngs[0])
	x01, y01 := proj.Forward(points.Lats[1], points.Lngs[1])
	x10, y10 := proj.Forward(points.Lats[ni], points.Lngs[ni])

	const tolerance = 1e-3 // meters
	assertMagnitude(t, "X change moving one column east", x01-x00, dxMeters, tolerance)
	assertMagnitude(t, "Y change moving one column east", y01-y00, 0, tolerance)
	assertMagnitude(t, "X change moving one row south", x10-x00, 0, tolerance)
	assertMagnitude(t, "Y change moving one row south", y10-y00, dyMeters, tolerance)
}
