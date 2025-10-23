package projections

import (
	"math"
)

type LambertConformalConicalParams struct {
	Radius                 float64
	Eccentricity           float64
	OriginLatitude         float64
	OriginLongitude        float64
	FirstStandardParallel  float64
	SecondStandardParallel float64
	Di                     float64 // must be in meters
	Dj                     float64 // must be in meters
	Ni                     int
	Nj                     int
}

func ExtractLambertConformalConicalGrid(params LambertConformalConicalParams) (lats []float32, lngs []float32) {
	originI := 0.0
	originJ := 0.0
	totalPoints := params.Ni * params.Nj
	p := newLambertConformalConic(
		params.Radius,
		params.Eccentricity,
		params.OriginLatitude,
		params.OriginLongitude,
		params.FirstStandardParallel,
		params.SecondStandardParallel,
		0.0,
		0.0,
	)
	lats, lngs = make([]float32, 0, totalPoints), make([]float32, 0, totalPoints)

	for j := 0; j < params.Nj; j++ {
		y := originJ + float64(j)*params.Dj
		for i := 0; i < params.Ni; i++ {
			x := originI + float64(i)*params.Di
			lat, lng := p.Inverse(x, y)
			lats = append(lats, float32(lat))
			lngs = append(lngs, float32(lng))
		}
	}
	return lats, lngs
}

// lambertConformalConic holds the projection parameters
type lambertConformalConic struct {
	// Ellipsoid parameters
	a  float64 // Semi-major axis (equatorial radius)
	e  float64 // Eccentricity
	e2 float64 // Eccentricity squared

	// Projection parameters
	lat0 float64 // Latitude of origin (radians)
	lon0 float64 // Longitude of origin (radians)
	lat1 float64 // First standard parallel (radians)
	lat2 float64 // Second standard parallel (radians)

	// False easting/northing
	falseEasting  float64
	falseNorthing float64

	// Computed constants
	n    float64 // Cone constant
	F    float64 // Scaling constant
	rho0 float64 // Radius at origin
}

// NewLambertConformalConic creates a new Lambert Conformal Conic projection
func newLambertConformalConic(radius, eccentricity, lat0Deg, lon0Deg, lat1Deg, lat2Deg, falseEasting, falseNorthing float64) *lambertConformalConic {
	lcc := &lambertConformalConic{
		a:  radius,                      // meters
		e:  eccentricity,                // eccentricity
		e2: eccentricity * eccentricity, // eccentricity squared

		lat0: degToRad(lat0Deg),
		lon0: degToRad(lon0Deg),
		lat1: degToRad(lat1Deg),
		lat2: degToRad(lat2Deg),

		falseEasting:  falseEasting,
		falseNorthing: falseNorthing,
	}

	// Compute projection constants
	lcc.computeConstants()

	return lcc
}

// computeConstants calculates the projection constants n, F, and rho0
func (lcc *lambertConformalConic) computeConstants() {
	// Calculate m values for the standard parallels
	m1 := lcc.computeM(lcc.lat1)
	m2 := lcc.computeM(lcc.lat2)

	// Calculate t values for the standard parallels and origin
	t1 := lcc.computeT(lcc.lat1)
	t2 := lcc.computeT(lcc.lat2)
	tF := lcc.computeT(lcc.lat0)

	// Cone constant (n)
	if math.Abs(lcc.lat1-lcc.lat2) < 1e-10 {
		// Single standard parallel case
		lcc.n = math.Sin(lcc.lat1)
	} else {
		// Two standard parallels
		lcc.n = (math.Log(m1) - math.Log(m2)) / (math.Log(t1) - math.Log(t2))
	}

	// Scaling constant (F)
	lcc.F = m1 / (lcc.n * math.Pow(t1, lcc.n))

	// Radius at origin (rho0)
	lcc.rho0 = lcc.a * lcc.F * math.Pow(tF, lcc.n)
}

// computeM calculates the m value for a given latitude
// m = cos(lat) / sqrt(1 - e^2 * sin^2(lat))
func (lcc *lambertConformalConic) computeM(lat float64) float64 {
	sinLat := math.Sin(lat)
	return math.Cos(lat) / math.Sqrt(1-lcc.e2*sinLat*sinLat)
}

// computeT calculates the t value for a given latitude
// t = tan(π/4 - lat/2) / [(1 - e*sin(lat)) / (1 + e*sin(lat))]^(e/2)
func (lcc *lambertConformalConic) computeT(lat float64) float64 {
	sinLat := math.Sin(lat)
	esinLat := lcc.e * sinLat

	tanPart := math.Tan(math.Pi/4 - lat/2)
	conformalPart := math.Pow((1-esinLat)/(1+esinLat), lcc.e/2)

	return tanPart / conformalPart
}

// Inverse converts projected coordinates (X, Y) to geographic coordinates (latitude, longitude)
// X and Y should be in meters
// Returns latitude and longitude in degrees
func (lcc *lambertConformalConic) Inverse(x, y float64) (latDeg, lonDeg float64) {
	// Remove false easting and northing
	x -= lcc.falseEasting
	y -= lcc.falseNorthing

	// Calculate rho' (radius from cone apex)
	// rho' = ±sqrt((x)^2 + (rho0 - y)^2), taking the sign of n
	rhoPrime := math.Sqrt(math.Pow(x, 2) + math.Pow(lcc.rho0-y, 2))
	if lcc.n < 0 {
		rhoPrime = -rhoPrime
	}

	// Calculate theta' (angle from central meridian)
	// theta' = atan2(x, rho0 - y)
	thetaPrime := math.Atan2(x, lcc.rho0-y)

	// Calculate t'
	// t' = (rho' / (a * F))^(1/n)
	tPrime := math.Pow(rhoPrime/(lcc.a*lcc.F), 1/lcc.n)

	// Calculate longitude
	// lon = theta' / n + lon0
	lon := thetaPrime/lcc.n + lcc.lon0

	// Calculate latitude using iterative method
	// This accounts for the ellipsoidal Earth
	lat := lcc.computeLatFromT(tPrime)

	// Convert to degrees
	latDeg = radToDeg(lat)
	lonDeg = radToDeg(lon)

	return latDeg, lonDeg
}

// computeLatFromT calculates latitude from t using iterative method
// lat = π/2 - 2*arctan(t * [(1 - e*sin(lat)) / (1 + e*sin(lat))]^(e/2))
func (lcc *lambertConformalConic) computeLatFromT(t float64) float64 {
	// Initial estimate
	lat := math.Pi/2 - 2*math.Atan(t)

	// Iterate to refine (usually converges in 3-5 iterations)
	for i := 0; i < 10; i++ {
		sinLat := math.Sin(lat)
		esinLat := lcc.e * sinLat

		conformalPart := math.Pow((1-esinLat)/(1+esinLat), lcc.e/2)
		latNew := math.Pi/2 - 2*math.Atan(t*conformalPart)

		// Check convergence (1e-10 radians ≈ 6e-9 degrees ≈ 0.0006mm at equator)
		if math.Abs(latNew-lat) < 1e-10 {
			return latNew
		}
		lat = latNew
	}

	return lat
}
