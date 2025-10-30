package record

import (
	"fmt"

	u "github.com/skysparq/grib2-go/utility"
)

// Section4 contains the fields from section 4 of a GRIB record.
// The product definition template is stored as raw bytes to defer processing until the user is ready.
type Section4 struct {
	Length                          int
	CoordinateValuesAfterTemplate   int
	ProductDefinitionTemplateNumber int
	ProductDefinitionTemplateData   []byte
	CoordinateValuesData            []byte
	Templates                       Templates
}

// ParseSection4 parses section 4 of a GRIB record.
// No attempt is made to parse the product definition template during this process.
// If the user needs to parse the product definition, call Section4.Definition after parsing the section.
// This allows the user to decide when to parse the product definition, and also allows the user
// to parse grib2 records that are not currently supported.
func ParseSection4(data SectionData, templates Templates) (section Section4, err error) {
	section.Length = data.Length
	if data.SectionNumber != 4 {
		return section, fmt.Errorf(`error parsing section 4: expected section number 4, got %d`, data.SectionNumber)
	}
	section.CoordinateValuesAfterTemplate = u.Uint16(data.Bytes[5:7])
	section.ProductDefinitionTemplateNumber = u.Uint16(data.Bytes[7:9])

	templateEnd, ok := templates.ProductDefinitionEnd(section.ProductDefinitionTemplateNumber, data.Bytes)
	if !ok {
		return section, fmt.Errorf(`error parsing section 4: unsupported Product} Definition Templates %d`, section.ProductDefinitionTemplateNumber)
	}
	if templateEnd-4 > data.Length {
		return section, fmt.Errorf(`error parsing section 4: Templates ending position %d exceeds available length %d`, templateEnd, data.Length)
	}

	section.ProductDefinitionTemplateData = data.Bytes[9:templateEnd]
	section.CoordinateValuesData = data.Bytes[templateEnd:]
	section.Templates = templates
	return section, nil
}

// Definition parses the product definition template and returns a ProductDefinition from the provided templates.
func (s Section4) Definition() (ProductDefinition, error) {
	return s.Templates.ProductDefinition(s)
}
