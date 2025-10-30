package projections

import u "github.com/skysparq/grib2-go/utility"

// EquidistantCylindricalParams contains the parameters needed to instantiate the Equidistant Cylindrical projection.
// This is a simple projection that defines a constant angular distance between points. For example, GFS gribs use this
// projection to generate points that are 1/4 degree of longitude and latitude apart. It is a constant-angle projection
// rather than a constant-distance projection, meaning spacing between points becomes smaller as the latitude approaches
// the poles.
type EquidistantCylindricalParams struct {
	ScanningMode ScanningMode
	Ni           int // total number of points in the x direction
	Nj           int // total number of points in the y direction
	Di           int // x increment (distance between points)
	Dj           int // y increment (distance between points)
	I0           int // starting x point
	J0           int // starting y point
}

// ExtractEquidistantCylindricalGrid extracts the grid points from the Equidistant Cylindrical projection defined by the given EquidistantCylindricalParams.
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
