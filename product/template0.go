package product

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template0 contains the fields for analysis or forecast at a horizontal level or in a horizontal layer at a point in time
type Template0 struct {
	ParameterCategory           int
	ParameterNumber             int
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

// Header returns the standard header fields common to all products
func (t Template0) Header(info record.Section1) record.ProductDefinitionHeader {
	at := u.TimestampFromReference(info.Time(), t.ForecastTimeInUnits, t.UnitOfTimeRange)
	return record.ProductDefinitionHeader{
		ParameterCategory:  t.ParameterCategory,
		ParameterNumber:    t.ParameterNumber,
		FirstSurfaceType:   t.FirstSurfaceType,
		FirstSurfaceValue:  u.ScaleInt(t.FirstSurfaceScaleValue, t.FirstSurfaceScaleFactor),
		SecondSurfaceType:  t.SecondSurfaceType,
		SecondSurfaceValue: u.ScaleInt(t.SecondSurfaceScaleValue, t.SecondSurfaceScaleFactor),
		Start:              at,
		End:                at,
		TimeIncrements:     nil,
	}
}

// Parse fills in the template from the provided section
func (t Template0) Parse(section record.Section4) (record.ProductDefinition, error) {
	err := checkSectionNum(section, 0)
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
