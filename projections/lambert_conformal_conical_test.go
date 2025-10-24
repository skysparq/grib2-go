package projections

import (
	"math"
	"testing"
)

func TestLambertConcormalConical(t *testing.T) {
	params := LambertConformalConicalParams{
		ScanningMode: ScanningMode{
			RightToLeft: false,
			TopToBottom: false,
			OverFirst:   true,
		},
		Radius:                 6371229.0,
		Eccentricity:           0.0,
		OriginLatitude:         38.5,
		OriginLongitude:        262.5 - 360.0,
		FirstStandardParallel:  38.5,
		SecondStandardParallel: 38.5,
		Di:                     3000,
		Dj:                     3000,
		Ni:                     1799,
		Nj:                     1059,
		StartLatitude:          21.138123,
		StartLongitude:         237.280472 - 360,
	}
	lats, lngs := ExtractLambertConformalConicalGrid(params)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			t.Logf("y: %d, x: %d, Lat: %f, Lng: %f\n", y, x, lats[y*1799+x], lngs[y*1799+x])
		}
	}
}

func TestLambertConcormalConical2(t *testing.T) {
	params := LambertConformalConicalParams{
		ScanningMode: ScanningMode{
			RightToLeft: false,
			TopToBottom: false,
			OverFirst:   true,
		},
		Radius:                 6371229.0,
		Eccentricity:           0.0,
		OriginLatitude:         38.5,
		OriginLongitude:        262.5 - 360.0,
		FirstStandardParallel:  38.5,
		SecondStandardParallel: 38.5,
	}
	l := newLambertConformalConic(params.Radius, params.Eccentricity, params.OriginLatitude, params.OriginLongitude, params.FirstStandardParallel, params.SecondStandardParallel, 0, 0)
	x, y := l.Forward(21.138123, 237.280472-360)
	t.Logf("X: %f, Y: %f\n", x, y)

	lat, lng := l.Inverse(x, y)
	t.Logf("Lat: %f, Lng: %f\n", lat, lng)
}

func TestLambertConcormalConical3(t *testing.T) {
	params := LambertConformalConicalParams{
		ScanningMode: ScanningMode{
			RightToLeft: false,
			TopToBottom: false,
			OverFirst:   true,
		},
		Radius:                 6378206.4,
		Eccentricity:           0.0822719,
		OriginLatitude:         23.0,
		OriginLongitude:        -96.0,
		FirstStandardParallel:  33.0,
		SecondStandardParallel: 45.0,
	}
	l := newLambertConformalConic(params.Radius, params.Eccentricity, params.OriginLatitude, params.OriginLongitude, params.FirstStandardParallel, params.SecondStandardParallel, 0, 0)
	lat, lng := l.Inverse(1894410.9, 1564649.5)
	t.Logf("Lat: %f, Lng: %f\n", lat, lng)
	if math.Abs(lat-35.0) > 1e-5 {
		t.Fatalf("Expected 35 but got %f", lat)
	}
	if math.Abs(lng - -75.0) > 1e-5 {
		t.Fatalf("Expected -75 but got %f", lng)
	}
}
