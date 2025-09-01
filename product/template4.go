package product

import (
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template4 struct {
	DefinitionHeader
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
	CentralPointLatitude         int
	CentralPointLongitude        int
	ClusterRadius                int
	TotalForecastsInCluster      int
	ScaledFactorStdDev           int
	ScaledValueStdDev            int
	ScaledFactorMean             int
	ScaledValueMean              int
	EnsembleForecastNumbers      []int
}

func (t Template4) Header() DefinitionHeader {
	return t.DefinitionHeader
}

func (t Template4) Parse(section record.Section4) (Definition, error) {
	err := checkSectionNum(section, 4)
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
	t.CentralPointLatitude = u.Int32(data[32:36])
	t.CentralPointLongitude = u.Int32(data[36:40])
	t.ClusterRadius = u.Int32(data[40:44])
	t.TotalForecastsInCluster = int(data[44])
	t.ScaledFactorStdDev = int(data[45])
	t.ScaledValueStdDev = u.Int32(data[46:50])
	t.ScaledFactorMean = int(data[50])
	t.ScaledValueMean = u.Int32(data[51:55])
	forecastNums := make([]int, len(data)-55)
	for i := range forecastNums {
		forecastNums[i] = int(data[i+55])
	}
	t.EnsembleForecastNumbers = forecastNums
	return t, nil
}
