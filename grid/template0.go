package grid

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

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

func (t Template0) Parse(section record.Section3) (Definition, error) {
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
