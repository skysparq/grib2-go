package record

import (
	"encoding/binary"
	"fmt"
	"github.com/skysparq/grib2-go/templates"
)

type Section3 struct {
	Length                       int
	GridSourceDefinition         int
	TotalPoints                  int
	OctetsForOptionalPointList   int
	InterpretationOfPointList    int
	GridDefinitionTemplateNumber int
	GridDefinitionTemplateData   []byte
	OptionalPointListData        []byte
}

func ParseSection3(data SectionData, template templates.Template) (section Section3, err error) {
	section.Length = data.Length
	if data.SectionNumber != 3 {
		return section, fmt.Errorf(`error parsing section 3: expected section number 3, got %d`, data.SectionNumber)
	}
	section.GridSourceDefinition = int(data.Bytes[5])
	section.TotalPoints = int(binary.BigEndian.Uint32(data.Bytes[6:10]))
	section.OctetsForOptionalPointList = int(data.Bytes[10])
	section.InterpretationOfPointList = int(data.Bytes[11])
	section.GridDefinitionTemplateNumber = int(binary.BigEndian.Uint16(data.Bytes[12:14]))
	templateEnd, ok := template.GridDefinitionEnd(section.GridDefinitionTemplateNumber, data.Bytes)
	if !ok {
		return section, fmt.Errorf(`error parsing section 3: unsupported Grid Definition Template %d`, section.GridDefinitionTemplateNumber)
	}
	if templateEnd-4 > data.Length {
		return section, fmt.Errorf(`error parsing section 3: template ending position %d exceeds available length %d`, templateEnd, data.Length)
	}

	section.GridDefinitionTemplateData = data.Bytes[14:templateEnd]
	section.OptionalPointListData = data.Bytes[templateEnd:]
	return section, nil
}
