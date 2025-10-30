package record

import (
	"fmt"

	u "github.com/skysparq/grib2-go/utility"
)

// Section5 contains the fields from section 5 of a GRIB record.
// The data representation definition template is stored as raw bytes to defer processing until the user is ready.
type Section5 struct {
	Length                           int
	TotalDataPoints                  int
	DataRepresentationTemplateNumber int
	DataRepresentationTemplateData   []byte
	Templates                        Templates
}

// ParseSection5 parses section 5 of a GRIB record.
// No attempt is made to parse the data representation definition template during this process.
// If the user needs to parse the data representation definition, call Section5.Definition after parsing the section.
// This allows the user to decide when to parse the data representation definition, and also allows the user
// to parse grib2 records that are not currently supported.
func ParseSection5(data SectionData, templates Templates) (section Section5, err error) {
	section.Length = data.Length
	if data.SectionNumber != 5 {
		return section, fmt.Errorf(`error parsing section 5: expected section number 5, got %d`, data.SectionNumber)
	}
	section.TotalDataPoints = u.Uint32(data.Bytes[5:9])
	section.DataRepresentationTemplateNumber = u.Uint16(data.Bytes[9:11])
	section.DataRepresentationTemplateData = data.Bytes[11:]
	section.Templates = templates
	return section, nil
}

// Definition parses the data representation definition template and returns a DataRepresentationDefinition from the provided templates.
func (s Section5) Definition() (DataRepresentationDefinition, error) {
	return s.Templates.DataRepresentation(s)
}
