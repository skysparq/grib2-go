package product

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template2 contains the fields for derived forecasts based on all ensemble members at a horizontal level or in a horizontal layer at a point in time.
type Template2 struct {
	record.ProductDefinitionHeader
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
	DerivedForecast             int
	TotalForecastsInEnsemble    int
}

// Header returns the standard header fields common to all products
func (t Template2) Header() record.ProductDefinitionHeader {
	return t.ProductDefinitionHeader
}

// Parse fills in the template from the provided section
func (t Template2) Parse(section record.Section4) (record.ProductDefinition, error) {
	err := checkSectionNum(section, 2)
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
	t.DerivedForecast = int(data[25])
	t.TotalForecastsInEnsemble = int(data[26])
	return t, nil
}
