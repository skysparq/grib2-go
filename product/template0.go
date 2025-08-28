package product

import (
	"fmt"
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template0 struct {
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

func (t *Template0) Header() DefinitionHeader {
	return t.DefinitionHeader
}

func (t *Template0) Parse(section record.Section4) error {
	if section.ProductDefinitionTemplateNumber != 0 {
		return fmt.Errorf(`error parsing product definition template 0: section 4 template number is %d rather than 0`, section.ProductDefinitionTemplateNumber)
	}
	t.ParameterCategory = int(section.ProductDefinitionTemplateData[0])
	t.ParameterNumber = int(section.ProductDefinitionTemplateData[1])
	t.GeneratingProcessType = int(section.ProductDefinitionTemplateData[2])
	t.BackgroundIdentifier = int(section.ProductDefinitionTemplateData[3])
	t.GeneratingProcessIdentifier = int(section.ProductDefinitionTemplateData[4])
	t.HoursAfterReference = u.Uint16(section.ProductDefinitionTemplateData[5:7])
	t.MinutesAfterReference = int(section.ProductDefinitionTemplateData[7])
	t.UnitOfTimeRange = int(section.ProductDefinitionTemplateData[8])
	t.ForecastTimeInUnits = u.Uint32(section.ProductDefinitionTemplateData[9:13])
	t.FirstSurfaceType = int(section.ProductDefinitionTemplateData[13])
	t.FirstSurfaceScaleFactor = int(section.ProductDefinitionTemplateData[14])
	t.FirstSurfaceScaleValue = u.Int32(section.ProductDefinitionTemplateData[15:19])
	t.SecondSurfaceType = int(section.ProductDefinitionTemplateData[19])
	t.SecondSurfaceScaleFactor = int(section.ProductDefinitionTemplateData[20])
	t.SecondSurfaceScaleValue = u.Int32(section.ProductDefinitionTemplateData[21:25])
	return nil
}
