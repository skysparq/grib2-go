package grid

import (
	"fmt"
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

func (t *Template0) Parse(section record.Section3) error {
	if section.GridDefinitionTemplateNumber != 0 {
		return fmt.Errorf(`error parsing grid definition template 0: section 3 template number is %d rather than 0`, section.GridDefinitionTemplateNumber)
	}
	t.EarthShape = int(section.GridDefinitionTemplateData[0])
	t.RadiusScaleFactor = int(section.GridDefinitionTemplateData[1])
	t.RadiusScaleValue = u.Int32(section.GridDefinitionTemplateData[2:6])
	t.MajorAxisScaleFactor = int(section.GridDefinitionTemplateData[6])
	t.MajorAxisScaleValue = u.Int32(section.GridDefinitionTemplateData[7:11])
	t.MinorAxisScaleValue = int(section.GridDefinitionTemplateData[11])
	t.MinorAxisScaleFactor = u.Int32(section.GridDefinitionTemplateData[12:16])
	t.PointsAlongParallel = u.Int32(section.GridDefinitionTemplateData[16:20])
	t.PointsAlongMeridian = u.Int32(section.GridDefinitionTemplateData[20:24])
	t.BasicAngle = u.Int32(section.GridDefinitionTemplateData[24:28])
	t.Subdivisions = u.Int32(section.GridDefinitionTemplateData[28:32])
	t.FirstLatitude = u.SignAndMagnitudeInt32(section.GridDefinitionTemplateData[32:36])
	t.FirstLongitude = u.SignAndMagnitudeInt32(section.GridDefinitionTemplateData[36:40])
	t.ResolutionAndComponentFlags = section.GridDefinitionTemplateData[40]
	t.LastLatitude = u.SignAndMagnitudeInt32(section.GridDefinitionTemplateData[41:45])
	t.LastLongitude = u.SignAndMagnitudeInt32(section.GridDefinitionTemplateData[45:49])
	t.ParallelIncrement = u.Int32(section.GridDefinitionTemplateData[49:53])
	t.MeridianIncrement = u.Int32(section.GridDefinitionTemplateData[53:57])
	t.ScanningMode = section.GridDefinitionTemplateData[57]
	return nil
}
