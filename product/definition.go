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

func StandardTemplates() map[int]Definition {
	return map[int]Definition{
		0: &Template0{},
		8: &Template8{},
	}
}

type Parser struct {
	Templates map[int]Definition
}

func (p *Parser) ParseDefinition(section record.Section4) (Definition, error) {
	if p.Templates == nil {
		p.Templates = StandardTemplates()
	}

	templateNum := section.ProductDefinitionTemplateNumber
	parser, ok := p.Templates[templateNum]
	if !ok {
		return nil, fmt.Errorf(`error parsing product definition: unsupported template number %d`, section.ProductDefinitionTemplateNumber)
	}
	return parser.Parse(section)
}
