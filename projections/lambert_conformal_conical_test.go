package projections_test

import (
	"testing"

	"github.com/skysparq/grib2-go/projections"
)

func TestLambertConcormalConical(t *testing.T) {
	params := projections.LambertConformalConicalParams{
		Radius:                 6371229.0,
		Eccentricity:           0,
		OriginLatitude:         21.138123,
		OriginLongitude:        237.280472 - 360.0,
		FirstStandardParallel:  38.5,
		SecondStandardParallel: 38.5,
		Di:                     3000,
		Dj:                     3000,
		Ni:                     1799,
		Nj:                     1059,
	}
	lats, lngs := projections.ExtractLambertConformalConicalGrid(params)
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			t.Logf("Lat: %f, Lng: %f\n", lats[y*10+x], lngs[y*10+x])
		}
	}
}
