package product

import (
	"time"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type TimeIncrement struct {
	StatisticalProcess         int
	TimeIncrementType          int
	StatisticalUnitOfTimeRange int
	StatisticalLengthOfTime    int
	SuccessiveUnitOfTimeRange  int
	SuccessiveLengthOfTime     int
}

type Template8 struct {
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
	EndYear                     int
	EndMonth                    int
	EndDay                      int
	EndHour                     int
	EndMinute                   int
	EndSecond                   int
	TotalTimeRanges             int
	MissingDataValues           int
	TimeRanges                  []TimeIncrement
}

func (t Template8) Header() DefinitionHeader {
	return t.DefinitionHeader
}

func (t Template8) Parse(section record.Section4) (Definition, error) {
	err := checkSectionNum(section, 8)
	if err != nil {
		return t, err
	}

	var data = section.ProductDefinitionTemplateData
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
	t.EndYear = u.Uint16(data[25:27])
	t.EndMonth = int(data[27])
	t.EndDay = int(data[28])
	t.EndHour = int(data[29])
	t.EndMinute = int(data[30])
	t.EndSecond = int(data[31])
	t.TotalTimeRanges = int(data[32])
	t.MissingDataValues = u.Uint32(data[33:37])

	t.TimeRanges = make([]TimeIncrement, 0, t.TotalTimeRanges)
	for startOctet := 37; startOctet < len(data); startOctet += 12 {
		t.TimeRanges = append(t.TimeRanges, TimeIncrement{
			StatisticalProcess:         int(data[startOctet]),
			TimeIncrementType:          int(data[startOctet+1]),
			StatisticalUnitOfTimeRange: int(data[startOctet+2]),
			StatisticalLengthOfTime:    u.Int32(data[startOctet+3 : startOctet+7]),
			SuccessiveUnitOfTimeRange:  int(data[startOctet+7]),
			SuccessiveLengthOfTime:     u.Uint32(data[startOctet+8 : startOctet+12]),
		})
	}
	return t, nil
}

func (t Template8) EndTime() time.Time {
	return time.Date(t.EndYear, time.Month(t.EndMonth), t.EndDay, t.EndHour, t.EndMinute, t.EndSecond, 0, time.UTC)
}
