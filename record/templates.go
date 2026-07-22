package record

import (
	"time"

	"github.com/skysparq/grib2-go/projections"
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
	// GetValues returns the decoded data value of every point in the grid, normalized to row-major
	// order with north at row 0 and west at column 0 - the same order GridDefinition.Points()
	// returns, regardless of the scanning mode the source GRIB2 file actually used. Values[i]
	// always corresponds to GridPoints.Lats[i]/Lngs[i] for the same grid.
	GetValues(rec Record) ([]float64, error)
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
	// Points returns the latitude and longitude of every point in the grid, normalized to
	// row-major order with north at row 0 and west at column 0, regardless of the scanning mode
	// the source GRIB2 file actually used. See GridPoints.
	Points() (GridPoints, error)
	XVals() int
	YVals() int
	SrsWkt() (string, error)
	// ScanMode returns the scanning mode the source GRIB2 file used to store this grid's points on
	// disk. It describes the file's raw layout only - it does NOT describe the order Points() or
	// DataReader.GetValues() return, since both of those are always normalized to north-up,
	// west-first row-major order regardless of this value.
	ScanMode() projections.ScanningMode
}

// GridPoints is the standard struct containing latitude and longitude values from a projection.
// Points are in row-major order with north at row 0 and west at column 0, regardless of the
// scanning mode the source GRIB2 file actually used.
type GridPoints struct {
	Lats []float64
	Lngs []float64
}

// GriddedValues is the standard struct containing latitude, longitude, and data values from a GRIB2 record.
// All three are aligned and normalized to row-major order with north at row 0 and west at column 0;
// Values[i] corresponds to Lats[i]/Lngs[i].
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
