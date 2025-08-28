package product

import (
	"fmt"

	"github.com/skysparq/grib2/record"
)

type Definition interface {
	Header() DefinitionHeader
	Parse(section record.Section4) (Definition, error)
}

type DefinitionHeader struct {
	ParameterCategory int
	ParameterNumber   int
}

func ParseDefinition(section record.Section4) (Definition, error) {
	switch section.ProductDefinitionTemplateNumber {
	case 0:
		return (&Template0{}).Parse(section)
	case 8:
		return (&Template8{}).Parse(section)
	default:
		return nil, fmt.Errorf(`error parsing product definition: unsupported template number %d`, section.ProductDefinitionTemplateNumber)
	}
}
