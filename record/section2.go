package record

import (
	"fmt"
)

type Section2 struct {
	Length   int
	LocalUse []byte
}

func ParseSection2(data SectionData) (section Section2, err error) {
	section.Length = data.Length

	if data.SectionNumber != 2 {
		return section, fmt.Errorf(`error parsing section 2: expected section number 2, got %d`, data.SectionNumber)
	}
	section.LocalUse = data.Bytes[5:]
	return section, nil
}
