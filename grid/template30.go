package grid

import (
	"errors"
	"fmt"
	"math"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template30 struct {
	EarthShape                  int
	RadiusScaleFactor           int
	RadiusScaleValue            int
	MajorAxisScaleFactor        int
	MajorAxisScaleValue         int
	MinorAxisScaleFactor        int
	MinorAxisScaleValue         int
	Nx                          int
	Ny                          int
	La1                         int
	Lo1                         int
	ResolutionAndComponentFlags byte
	LaD                         int
	LoV                         int
	Dx                          int
	Dy                          int
	ProjectionCenterFlags       byte
	ScanningMode                byte
	Latin1                      int
	Latin2                      int
	SouthernPoleLatitude        int
	SouthernPoleLongitude       int
}

func (t Template30) Points() (record.GridPoints, error) {
	var result record.GridPoints
	var params lCCParams

	if t.MajorAxisScaleValue != 0 {
		return result, errors.New("error getting points: non-standard lat/lon scaling not implemented")
	}

	switch t.EarthShape {
	case 0:
		params.EarthRadius = 6367470.0
	case 6:
		params.EarthRadius = 6371229.0
	default:
		return result, errors.New("error getting points: unsupported earth shape")
	}

	params.Nx = t.Nx
	params.Ny = t.Ny
	params.Dx = float64(t.Dx) / 1000 // Assume millimeters, adjust as needed
	params.Dy = float64(t.Dy) / 1000 // Assume millimeters, adjust as needed
	params.Lat0 = float64(t.LaD) / 1000000
	params.Lon0 = float64(t.LoV) / 1000000
	params.Lat1 = float64(t.Latin1) / 1000000
	params.Lat2 = float64(t.Latin2) / 1000000
	params.X0 = 0.0 // Projection origin, not Lo1
	params.Y0 = 0.0 // Projection origin, not La1

	// Step 1: Compute cone constant (n)
	phi1 := degToRad(params.Lat1)
	phi2 := degToRad(params.Lat2)
	var n float64
	if math.Abs(phi1-phi2) < 1e-10 { // Tangent case
		n = math.Sin(phi1)
	} else { // Secant case
		n = (math.Log(math.Cos(phi1)) - math.Log(math.Cos(phi2))) /
			(math.Log(math.Tan(math.Pi/4+phi2/2)) - math.Log(math.Tan(math.Pi/4+phi1/2)))
	}

	// Step 2: Compute scaling factor (F)
	F := math.Cos(phi1) * math.Pow(math.Tan(math.Pi/4+phi1/2), n) / n

	for i := 0; i < t.Nx; i++ {
		for j := 0; j < t.Ny; j++ {
			// Step 3: Convert grid indices to (x, y) in projected plane
			x := params.X0 + float64(i)*params.Dx
			y := params.Y0 + float64(j)*params.Dy

			// Step 4: Compute polar coordinates (ρ, θ)
			xPrime := x - params.X0
			yPrime := y - params.Y0
			rho := math.Sqrt(xPrime*xPrime + yPrime*yPrime)
			var lat, lon float64
			if rho < 1e-10 { // Handle near-pole case
				lat = radToDeg(math.Pi / 2)
				lon = params.Lon0
			} else {
				theta := math.Atan2(xPrime, yPrime) / n
				lon = radToDeg(theta + degToRad(params.Lon0))
				lat = radToDeg(2*math.Atan(math.Pow(params.EarthRadius*F/rho, 1/n)) - math.Pi/2)
			}

			// Normalize longitude to -180–180°
			if lon > 180 {
				lon -= 360
			} else if lon < -180 {
				lon += 360
			}

			result.Lngs = append(result.Lngs, float32(lon))
			result.Lats = append(result.Lats, float32(lat))
		}
	}

	// Validate first grid point against La1/Lo1
	firstLat := float64(t.La1) / 1000000
	firstLon := float64(t.Lo1) / 1000000
	if firstLon > 180 {
		firstLon -= 360
	}
	if math.Abs(float64(result.Lats[0])-firstLat) > 0.01 || math.Abs(float64(result.Lngs[0])-firstLon) > 0.01 {
		fmt.Printf("Warning: First grid point (lat: %.6f, lon: %.6f) does not match expected (lat: %.6f, lon: %.6f)\n",
			result.Lats[0], result.Lngs[0], firstLat, firstLon)
	}

	return result, nil
}

func (t Template30) Parse(section record.Section3) (record.GridDefinition, error) {
	err := checkSectionNum(section, 30)
	if err != nil {
		return t, err
	}

	data := section.GridDefinitionTemplateData
	t.EarthShape = int(data[0])
	t.RadiusScaleFactor = int(data[1])
	t.RadiusScaleValue = u.Int32(data[2:6])
	t.MajorAxisScaleFactor = int(data[6])
	t.MajorAxisScaleValue = u.Int32(data[7:11])
	t.MinorAxisScaleValue = int(data[11])
	t.MinorAxisScaleFactor = u.Int32(data[12:16])
	t.Nx = u.Int32(data[16:20])
	t.Ny = u.Int32(data[20:24])
	t.La1 = u.Int32(data[24:28])
	t.Lo1 = u.Int32(data[28:32])
	t.ResolutionAndComponentFlags = data[32]
	t.LaD = u.Int32(data[33:37])
	t.LoV = u.Int32(data[37:41])
	t.Dx = u.Int32(data[41:45])
	t.Dy = u.Int32(data[45:49])
	t.ProjectionCenterFlags = data[49]
	t.ScanningMode = data[50]
	t.Latin1 = u.Int32(data[51:55])
	t.Latin2 = u.Int32(data[55:59])
	t.SouthernPoleLatitude = u.Int32(data[59:63])
	t.SouthernPoleLongitude = u.Int32(data[63:67])
	return t, nil
}
func degToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// radToDeg converts radians to degrees
func radToDeg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

type lCCParams struct {
	Nx, Ny      int     // Number of grid points
	Dx, Dy      float64 // Grid spacing in meters
	Lat0, Lon0  float64 // Reference latitude and longitude (degrees)
	Lat1, Lat2  float64 // Standard parallels (degrees)
	X0, Y0      float64 // Projection center coordinates (meters)
	EarthRadius float64 // Earth radius (meters)
}
