package record

import (
	"encoding/binary"
	"fmt"
)

type Section5 struct {
	Length                           int
	TotalDataPoints                  int
	DataRepresentationTemplateNumber int
	DataRepresentationTemplateData   []byte
}

func ParseSection5(data SectionData) (section Section5, err error) {
	section.Length = data.Length
	if data.SectionNumber != 5 {
		return section, fmt.Errorf(`error parsing section 5: expected section number 5, got %d`, data.SectionNumber)
	}
	section.TotalDataPoints = int(binary.BigEndian.Uint32(data.Bytes[5:9]))
	section.DataRepresentationTemplateNumber = int(binary.BigEndian.Uint16(data.Bytes[9:11]))
	section.DataRepresentationTemplateData = data.Bytes[11:]
	return section, nil
}
