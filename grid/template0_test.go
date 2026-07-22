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
		t.Fatalf("expected last latitude to be -90 but got %v", value)
	}

	if value := points.Lngs[0]; value != 0 {
		t.Fatalf("expected first longitude to be 0 but got %v", value)
	}

	if value := points.Lngs[1038239]; value != -0.25 {
		t.Fatalf("expected last longitude to be -0.25 but got %v", value)
	}
}

// TestTemplate0SrsWktMatchesPoints verifies that SrsWkt() actually describes the same coordinate
// system Points() used to generate the grid: it re-expresses the (already normalized) generated
// lat/lng points relative to the WKT's PRIMEM shift - the same reprojection a GEOGCS-to-GEOGCS
// transform through this WKT would perform - normalized into [0, 360) so the result stays
// monotonic across the antimeridian, and checks the result is an exactly regular Ni x Nj grid
// spaced by the grid's declared increments. That's only possible if the WKT's prime meridian
// shift matches the grid's own starting longitude.
func TestTemplate0SrsWktMatchesPoints(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordProdDef0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	definition, err := grid.Template0{}.Parse(rec.Grid)
	if err != nil {
		t.Fatal(err)
	}
	template := definition.(grid.Template0)

	points, err := template.Points()
	if err != nil {
		t.Fatal(err)
	}
	wkt, err := template.SrsWkt()
	if err != nil {
		t.Fatal(err)
	}

	shift := wktPrimeMeridian(t, wkt)
	forward := func(lat, lng float64) (x, y float64) {
		x = math.Mod(lng-shift, 360)
		if x < 0 {
			x += 360
		}
		return x, lat
	}

	// Points() is normalized to row-major order, so index 0 is the north-west corner, index 1 is
	// one column east, and index Ni is one row south.
	ni := template.XVals()
	dx := float64(template.ParallelIncrement) * 1e-6
	dy := float64(template.MeridianIncrement) * 1e-6

	x00, y00 := forward(points.Lats[0], points.Lngs[0])
	x01, y01 := forward(points.Lats[1], points.Lngs[1])
	x10, y10 := forward(points.Lats[ni], points.Lngs[ni])

	const tolerance = 1e-6 // degrees
	assertMagnitude(t, "X change moving one column east", x01-x00, dx, tolerance)
	assertMagnitude(t, "Y change moving one column east", y01-y00, 0, tolerance)
	assertMagnitude(t, "X change moving one row south", x10-x00, 0, tolerance)
	assertMagnitude(t, "Y change moving one row south", y10-y00, dy, tolerance)
}
