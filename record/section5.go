package record

import (
	"fmt"

	u "github.com/skysparq/grib2-go/utility"
)

type Section5 struct {
	Length                           int
	TotalDataPoints                  int
	DataRepresentationTemplateNumber int
	DataRepresentationTemplateData   []byte
	Templates                        Templates
}

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

func (s Section5) DataRepresentationDefinition() (DataRepresentationDefinition, error) {
	return s.Templates.DataRepresentation(s)
}
