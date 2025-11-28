package record

import (
	"iter"
	"time"
)

// Templates provides an interface for passing GRIB2 section templates to the parser.
type Templates interface {
	GridDefinitionEnd(template int, section3Bytes []byte) (int, bool)    // Table 3.1
	ProductDefinitionEnd(template int, section4Bytes []byte) (int, bool) // Table 4.0
	DataRepresentation(section Section5) (DataRepresentationDefinition, error)
	GridDefinition(section Section3) (GridDefinition, error)
	ProductDefinition(section Section4) (ProductDefinition, error)
}

// DataRepresentationDefinition provides an interface for parsing GRIB2 section 5.
// It includes methods for retrieving certain standardized information from the section.
type DataRepresentationDefinition interface {
	Parse(section Section5) (DataRepresentationDefinition, error)
	DecimalScale() int
	DataReader
}

// DataReader provides an interface for retrieving values from a GRIB2 record.
type DataReader interface {
	GetValues(rec Record) ([]float64, error)
	ValuesIterator(rec Record) (iter.Seq2[int, float64], error)
}

// ProductDefinition provides an interface for parsing GRIB2 section 4.
// It includes methods for retrieving certain standardized information from the section.
type ProductDefinition interface {
	Header(info Section1) ProductDefinitionHeader
	Parse(section Section4) (ProductDefinition, error)
}

// GridDefinition provides an interface for parsing GRIB2 section 3.
// It includes methods for retrieving certain standardized information from the section.
type GridDefinition interface {
	Parse(section Section3) (GridDefinition, error)
	Points() (GridPoints, error)
	XVals() int
	YVals() int
}

// GridPoints is the standard struct containing latitude and longitude values from a projection.
type GridPoints struct {
	Lats []float64
	Lngs []float64
}

// GriddedValues is the standard struct containing latitude, longitude, and data values from a GRIB2 record.
type GriddedValues struct {
	XVals int
	YVals int
	GridPoints
	Values []float64
}

// ProductDefinitionHeader contains standard fields present in every product definition section.
type ProductDefinitionHeader struct {
	ParameterCategory  int
	ParameterNumber    int
	FirstSurfaceType   int
	FirstSurfaceValue  float64
	SecondSurfaceType  int
	SecondSurfaceValue float64
	Start              time.Time
	End                time.Time
	TimeIncrements     []TimeIncrement
}

// TimeIncrement defines the time intervals used in certain product definition templates where statistical processing is performed..
type TimeIncrement struct {
	StatisticalProcess         int
	TimeIncrementType          int
	StatisticalUnitOfTimeRange int
	StatisticalLengthOfTime    int
	SuccessiveUnitOfTimeRange  int
	SuccessiveLengthOfTime     int
}
