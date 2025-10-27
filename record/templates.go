package record

type Templates interface {
	GridDefinitionEnd(template int, section3Bytes []byte) (int, bool)    // Table 3.1
	ProductDefinitionEnd(template int, section4Bytes []byte) (int, bool) // Table 4.0
	DataRepresentation(section Section5) (DataRepresentationDefinition, error)
	GridDefinition(section Section3) (GridDefinition, error)
	ProductDefinition(section Section4) (ProductDefinition, error)
}

type DataRepresentationDefinition interface {
	Parse(section Section5) (DataRepresentationDefinition, error)
	DataReader
}

type DataReader interface {
	GetValues(rec Record) ([]float64, error)
}

type ProductDefinition interface {
	Header() ProductDefinitionHeader
	Parse(section Section4) (ProductDefinition, error)
}

type GridDefinition interface {
	Parse(section Section3) (GridDefinition, error)
	Points() (GridPoints, error)
}

type GridPoints struct {
	Lats []float64
	Lngs []float64
}

type GriddedValues struct {
	GridPoints
	Values []float64
}

type ProductDefinitionHeader struct {
	ParameterCategory int
	ParameterNumber   int
}
