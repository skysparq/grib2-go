package product

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template3 struct {
	DefinitionHeader
	GeneratingProcessType       int
	BackgroundIdentifier        int
	GeneratingProcessIdentifier int
	HoursAfterReference         int
	MinutesAfterReference       int
	UnitOfTimeRange             int
	ForecastTimeInUnits         int
	FirstSurfaceType            int
	FirstSurfaceScaleFactor     int
	FirstSurfaceScaleValue      int
	SecondSurfaceType           int
	SecondSurfaceScaleFactor    int
	SecondSurfaceScaleValue     int
}

func (t Template3) Header() DefinitionHeader {
	return t.DefinitionHeader
}

func (t Template3) Parse(section record.Section4) (Definition, error) {
	err := checkSectionNum(section, 3)
	if err != nil {
		return t, err
	}

	data := section.ProductDefinitionTemplateData
	t.ParameterCategory = int(data[0])
	t.ParameterNumber = int(data[1])
	t.GeneratingProcessType = int(data[2])
	t.BackgroundIdentifier = int(data[3])
	t.GeneratingProcessIdentifier = int(data[4])
	t.HoursAfterReference = u.Uint16(data[5:7])
	t.MinutesAfterReference = int(data[7])
	t.UnitOfTimeRange = int(data[8])
	t.ForecastTimeInUnits = u.Uint32(data[9:13])
	t.FirstSurfaceType = int(data[13])
	t.FirstSurfaceScaleFactor = int(data[14])
	t.FirstSurfaceScaleValue = u.Int32(data[15:19])
	t.SecondSurfaceType = int(data[19])
	t.SecondSurfaceScaleFactor = int(data[20])
	t.SecondSurfaceScaleValue = u.Int32(data[21:25])
	return t, nil
}
