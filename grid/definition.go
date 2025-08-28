package grid

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
)

type Definition interface {
	Parse(section record.Section3) (Definition, error)
}

func StandardTemplates() map[int]Definition {
	return map[int]Definition{
		0: &Template0{},
	}
}

type Parser struct {
	Templates map[int]Definition
}

func (p *Parser) ParseDefinition(section record.Section3) (Definition, error) {
	if p.Templates == nil {
		p.Templates = StandardTemplates()
	}

	templateNum := section.GridDefinitionTemplateNumber
	parser, ok := p.Templates[templateNum]
	if !ok {
		return nil, fmt.Errorf(`error parsing grid definition: unsupported template number %d`, templateNum)
	}

	return parser.Parse(section)
}
