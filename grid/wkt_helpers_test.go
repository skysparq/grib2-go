package grid_test

import (
	"math"
	"regexp"
	"strconv"
	"testing"
)

// wktParam extracts the numeric value of a PARAMETER["name", value] entry from a WKT string.
func wktParam(t *testing.T, wkt, name string) float64 {
	t.Helper()
	re := regexp.MustCompile(`PARAMETER\["` + regexp.QuoteMeta(name) + `",\s*([-\d.]+)\]`)
	match := re.FindStringSubmatch(wkt)
	if match == nil {
		t.Fatalf("expected WKT to contain parameter %q, got: %s", name, wkt)
	}
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		t.Fatalf("error parsing parameter %q value %q: %v", name, match[1], err)
	}
	return value
}

// wktSpheroidRadius extracts the semi-major axis value from a WKT string's SPHEROID entry.
func wktSpheroidRadius(t *testing.T, wkt string) float64 {
	t.Helper()
	re := regexp.MustCompile(`SPHEROID\["[^"]*",\s*([-\d.]+),`)
	match := re.FindStringSubmatch(wkt)
	if match == nil {
		t.Fatalf("expected WKT to contain a SPHEROID definition, got: %s", wkt)
	}
	value, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		t.Fatalf("error parsing spheroid radius %q: %v", match[1], err)
	}
	return value
}

// assertMagnitude fails the test unless |got| is within tolerance of wantMagnitude. The sign of
// got is ignored since it depends on scan direction conventions that aren't under test here.
func assertMagnitude(t *testing.T, label string, got, wantMagnitude, tolerance float64) {
	t.Helper()
	if diff := math.Abs(math.Abs(got) - wantMagnitude); diff > tolerance {
		t.Fatalf("expected %s to have magnitude %v but got %v", label, wantMagnitude, got)
	}
}
