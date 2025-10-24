package projections_test

import (
	"testing"

	"github.com/skysparq/grib2-go/projections"
)

func TestEquidistantCylindricalTopLeftToBottomRight(t *testing.T) {
	params := projections.EquidistantCylindricalParams{
		ScanningMode: projections.ScanningMode{
			RightToLeft: false,
			TopToBottom: true,
			OverFirst:   true,
		},
		Ni: 10,
		Nj: 10,
		Di: 1000000,
		Dj: 2000000,
		I0: -170000000,
		J0: 20000000,
	}
	lats, lngs := projections.ExtractEquidistantCylindricalGrid(params)

	index := 0
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			expectedLon := float64(-170 + (j * 1))
			expectedLat := float64(20 - (i * 2))
			actualLon := lngs[index]
			actualLat := lats[index]
			if actualLat != expectedLat {
				t.Fatalf("expected latitude %v but got %v", expectedLat, actualLat)
			}
			if actualLon != expectedLon {
				t.Fatalf("expected longitude %v but got %v", expectedLon, actualLon)
			}
			index++
		}
	}
}
