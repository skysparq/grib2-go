package record

import (
	"encoding/binary"
	"fmt"
	"github.com/skysparq/grib2-go/templates"
)

type Section4 struct {
	Length                          int
	CoordinateValuesAfterTemplate   int
	ProductDefinitionTemplateNumber int
	ProductDefinitionTemplateData   []byte
	CoordinateValuesData            []byte
}

func ParseSection4(data SectionData, template templates.Template) (section Section4, err error) {
	section.Length = data.Length
	if data.SectionNumber != 4 {
		return section, fmt.Errorf(`error parsing section 4: expected section number 4, got %d`, data.SectionNumber)
	}
	section.CoordinateValuesAfterTemplate = int(binary.BigEndian.Uint16(data.Bytes[5:7]))
	section.ProductDefinitionTemplateNumber = int(binary.BigEndian.Uint16(data.Bytes[7:9]))

	templateEnd, ok := template.ProductDefinitionEnd(section.ProductDefinitionTemplateNumber, data.Bytes)
	if !ok {
		return section, fmt.Errorf(`error parsing section 4: unsupported Product} Definition Template %d`, section.ProductDefinitionTemplateNumber)
	}
	if templateEnd-4 > data.Length {
		return section, fmt.Errorf(`error parsing section 4: template ending position %d exceeds available length %d`, templateEnd, data.Length)
	}

	section.ProductDefinitionTemplateData = data.Bytes[9:templateEnd]
	section.CoordinateValuesData = data.Bytes[templateEnd:]

	return section, nil
}
