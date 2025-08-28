package product

import "github.com/skysparq/grib2-go/record"

type Definition interface {
	Header() DefinitionHeader
	Parse(section record.Section4) error
}

type DefinitionHeader struct {
	ParameterCategory int
	ParameterNumber   int
}
