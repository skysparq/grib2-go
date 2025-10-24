package projections

import u "github.com/skysparq/grib2-go/utility"

type EquidistantCylindricalParams struct {
	ScanningMode ScanningMode
	Ni           int // total number of points in the x direction
	Nj           int // total number of points in the y direction
	Di           int // x increment (distance between points)
	Dj           int // y increment (distance between points)
	I0           int // starting x point
	J0           int // starting y point
}

func ExtractEquidistantCylindricalGrid(params EquidistantCylindricalParams) (lats []float64, lngs []float64) {
	scannerParams := ScannerParams[int]{
		ScanningMode: params.ScanningMode,
		Ni:           params.Ni,
		Nj:           params.Nj,
		Di:           params.Di,
		Dj:           params.Dj,
		I0:           params.I0,
		J0:           params.J0,
	}

	s := NewScanner(scannerParams)
	lats = make([]float64, 0, params.Ni*params.Nj)
	lngs = make([]float64, 0, params.Ni*params.Nj)
	for lat, lng := range s.Points {
		lats = append(lats, u.StdLatLngToFloat(lat))
		lngs = append(lngs, u.StdLatLngToFloat(u.ShiftLongitude(lng)))
	}

	return lats, lngs
}
