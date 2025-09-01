package data_representation

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
)

type ValueGetter func(rec record.Record) ([]float32, error)

type Definition interface {
	Parse(section record.Section5) (Definition, error)
	DataReader
}

type DataReader interface {
	GetValues(rec record.Record) ([]float32, error)
}

func StandardTemplates() map[int]Definition {
	return map[int]Definition{
		0: Template0{},
		3: Template3{},
	}
}

type Parser struct {
	Templates map[int]Definition
}

func (p *Parser) ParseDefinition(section record.Section5) (Definition, error) {
	if p.Templates == nil {
		p.Templates = StandardTemplates()
	}

	templateNum := section.DataRepresentationTemplateNumber
	parser, ok := p.Templates[templateNum]
	if !ok {
		return nil, fmt.Errorf(`error parsing data representation: unsupported template number %d`, templateNum)
	}

	return parser.Parse(section)
}
