package grid

import (
	"errors"

	"github.com/skysparq/grib2-go/projections"
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template0 contains the fields for Latitude/Longitude
type Template0 struct {
	EarthShape                  int
	RadiusScaleFactor           int
	RadiusScaleValue            int
	MajorAxisScaleFactor        int
	MajorAxisScaleValue         int
	MinorAxisScaleFactor        int
	MinorAxisScaleValue         int
	PointsAlongParallel         int
	PointsAlongMeridian         int
	BasicAngle                  int
	Subdivisions                int
	FirstLatitude               int
	FirstLongitude              int
	ResolutionAndComponentFlags byte
	LastLatitude                int
	LastLongitude               int
	ParallelIncrement           int
	MeridianIncrement           int
	ScanningMode                byte
}

// Points returns the latitude and longitude for each point in the grid.
func (t Template0) Points() (record.GridPoints, error) {
	var result record.GridPoints
	if t.MajorAxisScaleValue != 0 {
		return result, errors.New("error getting points: non-standard lat/lon scaling not implemented")
	}

	params := projections.EquidistantCylindricalParams{
		ScanningMode: projections.ScanningModeFromByte(t.ScanningMode),
		Ni:           t.PointsAlongParallel,
		Nj:           t.PointsAlongMeridian,
		Di:           t.ParallelIncrement,
		Dj:           t.MeridianIncrement,
		I0:           u.ShiftLongitude(t.FirstLongitude),
		J0:           t.FirstLatitude,
	}
	result.Lats, result.Lngs = projections.ExtractEquidistantCylindricalGrid(params)

	return result, nil
}

// Parse fills in the template from the provided section
func (t Template0) Parse(section record.Section3) (record.GridDefinition, error) {
	err := checkSectionNum(section, 0)
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
	t.PointsAlongParallel = u.Int32(data[16:20])
	t.PointsAlongMeridian = u.Int32(data[20:24])
	t.BasicAngle = u.Int32(data[24:28])
	t.Subdivisions = u.Int32(data[28:32])
	t.FirstLatitude = u.SignAndMagnitudeInt32(data[32:36])
	t.FirstLongitude = u.SignAndMagnitudeInt32(data[36:40])
	t.ResolutionAndComponentFlags = data[40]
	t.LastLatitude = u.SignAndMagnitudeInt32(data[41:45])
	t.LastLongitude = u.SignAndMagnitudeInt32(data[45:49])
	t.ParallelIncrement = u.Int32(data[49:53])
	t.MeridianIncrement = u.Int32(data[53:57])
	t.ScanningMode = data[57]
	return t, nil
}

// XVals returns the number of points along the X axis.
func (t Template0) XVals() int {
	return t.PointsAlongParallel
}

// YVals returns the number of points along the Y axis.
func (t Template0) YVals() int {
	return t.PointsAlongMeridian
}
