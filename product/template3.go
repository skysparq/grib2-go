package product

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template3 contains the fields for derived forecasts based on a cluster of ensemble members over a rectangular area at a horizontal level or in a horizontal layer at a point in time.
type Template3 struct {
	ParameterCategory            int
	ParameterNumber              int
	GeneratingProcessType        int
	BackgroundIdentifier         int
	GeneratingProcessIdentifier  int
	HoursAfterReference          int
	MinutesAfterReference        int
	UnitOfTimeRange              int
	ForecastTimeInUnits          int
	FirstSurfaceType             int
	FirstSurfaceScaleFactor      int
	FirstSurfaceScaleValue       int
	SecondSurfaceType            int
	SecondSurfaceScaleFactor     int
	SecondSurfaceScaleValue      int
	DerivedForecast              int
	TotalForecastsInEnsemble     int
	ClusterId                    int
	HighResolutionControlCluster int
	LowResolutionControlCluster  int
	TotalClusters                int
	ClusteringMethod             int
	NorthernLat                  int
	SouthernLat                  int
	EasternLng                   int
	WesternLng                   int
	TotalForecastsInCluster      int
	ScaledFactorStdDev           int
	ScaledValueStdDev            int
	ScaledFactorMean             int
	ScaledValueMean              int
	EnsembleForecastNumbers      []int
}

// Header returns the standard header fields common to all products
func (t Template3) Header() record.ProductDefinitionHeader {
	return record.ProductDefinitionHeader{
		ParameterCategory:  t.ParameterCategory,
		ParameterNumber:    t.ParameterNumber,
		FirstSurfaceType:   t.FirstSurfaceType,
		FirstSurfaceValue:  u.ScaleInt(t.FirstSurfaceScaleValue, t.FirstSurfaceScaleFactor),
		SecondSurfaceType:  t.SecondSurfaceType,
		SecondSurfaceValue: u.ScaleInt(t.SecondSurfaceScaleValue, t.SecondSurfaceScaleFactor),
		TimeIncrements:     nil,
	}
}

// Parse fills in the template from the provided section
func (t Template3) Parse(section record.Section4) (record.ProductDefinition, error) {
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
	t.DerivedForecast = int(data[25])
	t.TotalForecastsInEnsemble = int(data[26])
	t.ClusterId = int(data[27])
	t.HighResolutionControlCluster = int(data[28])
	t.LowResolutionControlCluster = int(data[29])
	t.TotalClusters = int(data[30])
	t.ClusteringMethod = int(data[31])
	t.NorthernLat = u.Int32(data[32:36])
	t.SouthernLat = u.Int32(data[36:40])
	t.EasternLng = u.Int32(data[40:44])
	t.WesternLng = u.Int32(data[44:48])
	t.TotalForecastsInCluster = int(data[48])
	t.ScaledFactorStdDev = int(data[49])
	t.ScaledValueStdDev = u.Int32(data[50:54])
	t.ScaledFactorMean = int(data[54])
	t.ScaledValueMean = u.Int32(data[55:59])
	forecastNums := make([]int, len(data)-59)
	for i := range forecastNums {
		forecastNums[i] = int(data[i+59])
	}
	t.EnsembleForecastNumbers = forecastNums
	return t, nil
}
