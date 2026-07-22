package grid

import (
	"errors"
	"fmt"

	"github.com/skysparq/grib2-go/projections"
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template30 contains the fields for Lambert Conformal
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

// Points returns the latitude and longitude for each point in the grid, normalized to row-major
// order with north at row 0 and west at column 0.
func (t Template30) Points() (record.GridPoints, error) {
	var result record.GridPoints
	if t.MajorAxisScaleValue != 0 {
		return result, errors.New("error getting points: non-standard lat/lon scaling not implemented")
	}

	mode := t.ScanMode()
	params := projections.LambertConformalConicalParams{
		ScanningMode:           mode,
		OriginLatitude:         u.StdLatLngToFloat(t.LaD),
		OriginLongitude:        u.StdLatLngToFloat(u.ShiftLongitude(t.LoV)),
		FirstStandardParallel:  u.StdLatLngToFloat(t.Latin1),
		SecondStandardParallel: u.StdLatLngToFloat(t.Latin2),
		Di:                     mmToMeters(t.Dx),
		Dj:                     mmToMeters(t.Dy),
		Ni:                     t.Nx,
		Nj:                     t.Ny,
		StartLatitude:          u.StdLatLngToFloat(t.La1),
		StartLongitude:         u.StdLatLngToFloat(u.ShiftLongitude(t.Lo1)),
	}

	radius, err := earthRadius(t.EarthShape)
	if err != nil {
		return result, fmt.Errorf("error getting points: %w", err)
	}
	params.Radius = radius

	result.Lats, result.Lngs = projections.ExtractLambertConformalConicalGrid(params)
	result.Lats = projections.NormalizeScanOrder(result.Lats, t.Nx, t.Ny, mode)
	result.Lngs = projections.NormalizeScanOrder(result.Lngs, t.Nx, t.Ny, mode)

	return result, nil
}

// ScanMode returns the scanning mode used to order this grid's points.
func (t Template30) ScanMode() projections.ScanningMode {
	return projections.ScanningModeFromByte(t.ScanningMode)
}

// Parse fills in the template from the provided section
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

// XVals returns the number of points along the X axis.
func (t Template30) XVals() int {
	return t.Nx
}

// YVals returns the number of points along the Y axis.
func (t Template30) YVals() int {
	return t.Ny
}

// SrsWkt returns the WKT string describing the Lambert Conformal Conic projection.
func (t Template30) SrsWkt() (string, error) {
	radius, err := earthRadius(t.EarthShape)
	if err != nil {
		return "", fmt.Errorf(`error generating SRS WKT: %w`, err)
	}

	return fmt.Sprintf(
		`PROJCS["unnamed", GEOGCS["unnamed", DATUM["unknown", SPHEROID["unnamed", %.1f, 0]], PRIMEM["Greenwich", 0], UNIT["degree", 0.0174532925199433]], PROJECTION["Lambert_Conformal_Conic_2SP"], PARAMETER["False_Easting", 0], PARAMETER["False_Northing", 0], PARAMETER["Central_Meridian", %v], PARAMETER["Standard_Parallel_1", %v], PARAMETER["Standard_Parallel_2", %v], PARAMETER["Latitude_Of_Origin", %v], UNIT["Metre", 1]]`,
		radius,
		u.StdLatLngToFloat(u.ShiftLongitude(t.LoV)),
		u.StdLatLngToFloat(t.Latin1),
		u.StdLatLngToFloat(t.Latin2),
		u.StdLatLngToFloat(t.LaD),
	), nil
}
