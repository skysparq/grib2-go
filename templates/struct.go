package templates

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
)

type templates struct {
	gridDefinitionEndingOctet    map[int]RetrieveEndOctet
	productDefinitionEndingOctet map[int]RetrieveEndOctet
	dataRepresentationTemplates  map[int]record.DataRepresentationDefinition
	gridDefinitionTemplates      map[int]record.GridDefinition
	productDefinitionTemplates   map[int]record.ProductDefinition
}

func (t *templates) GridDefinitionEnd(template int, section3Bytes []byte) (int, bool) {
	retriever, ok := t.gridDefinitionEndingOctet[template]
	if !ok {
		return 0, false
	}
	return retriever(section3Bytes), true
}

func (t *templates) ProductDefinitionEnd(template int, section4Bytes []byte) (int, bool) {
	retriever, ok := t.productDefinitionEndingOctet[template]
	if !ok {
		return 0, false
	}
	return retriever(section4Bytes), true
}

func (t *templates) DataRepresentation(section record.Section5) (record.DataRepresentationDefinition, error) {
	templateNum := section.DataRepresentationTemplateNumber
	parser, ok := t.dataRepresentationTemplates[templateNum]
	if !ok {
		return nil, fmt.Errorf(`error parsing data representation: unsupported templates number %d`, templateNum)
	}

	return parser.Parse(section)
}

func (t *templates) GridDefinition(section record.Section3) (record.GridDefinition, error) {
	templateNum := section.GridDefinitionTemplateNumber
	parser, ok := t.gridDefinitionTemplates[templateNum]
	if !ok {
		return nil, fmt.Errorf(`error parsing grid definition: unsupported templates number %d`, templateNum)
	}

	return parser.Parse(section)
}

func (t *templates) ProductDefinition(section record.Section4) (record.ProductDefinition, error) {
	templateNum := section.ProductDefinitionTemplateNumber
	parser, ok := t.productDefinitionTemplates[templateNum]
	if !ok {
		return nil, fmt.Errorf(`error parsing product definition: unsupported templates number %d`, templateNum)
	}

	return parser.Parse(section)
}

type RetrieveEndOctet func([]byte) int
