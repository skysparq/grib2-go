package grid

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template30 struct {
	EarthShape                      int
	RadiusScaleFactor               int
	RadiusScaleValue                int
	MajorAxisScaleFactor            int
	MajorAxisScaleValue             int
	MinorAxisScaleFactor            int
	MinorAxisScaleValue             int
	PointsAlongParallel             int
	PointsAlongMeridian             int
	FirstLatitude                   int
	FirstLongitude                  int
	ResolutionAndComponentFlags     byte
	LatitudeDxDySpecified           int
	LongitudeWhereLatitudeIncreases int
	ParallelGridLength              int
	MeridianGridLength              int
	ProjectionCenterFlags           byte
	ScanningMode                    byte
	SecantConeFirstLatitude         int
	SecantConeSecondLatitude        int
	SouthernPoleLatitude            int
	SouthernPoleLongitude           int
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
	t.PointsAlongParallel = u.Int32(data[16:20])
	t.PointsAlongMeridian = u.Int32(data[20:24])
	t.FirstLatitude = u.Int32(data[24:28])
	t.FirstLongitude = u.Int32(data[28:32])
	t.ResolutionAndComponentFlags = data[32]
	t.LatitudeDxDySpecified = u.Int32(data[33:37])
	t.LongitudeWhereLatitudeIncreases = u.Int32(data[37:41])
	t.ParallelGridLength = u.Int32(data[41:45])
	t.MeridianGridLength = u.Int32(data[45:49])
	t.ProjectionCenterFlags = data[49]
	t.ScanningMode = data[50]
	t.SecantConeFirstLatitude = u.Int32(data[51:55])
	t.SecantConeSecondLatitude = u.Int32(data[55:59])
	t.SouthernPoleLatitude = u.Int32(data[59:63])
	t.SouthernPoleLongitude = u.Int32(data[63:67])
	return t, nil
}
