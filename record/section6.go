package record

import (
	"fmt"
)

type Section6 struct {
	Length          int
	BitmapIndicator int
	BitmapData      []byte
}

func ParseSection6(data SectionData) (section Section6, err error) {
	section.Length = data.Length
	if data.SectionNumber != 6 {
		return section, fmt.Errorf(`error parsing section 6: expected section number 6, got %d`, data.SectionNumber)
	}
	section.BitmapIndicator = int(data.Bytes[5])
	section.BitmapData = data.Bytes[6:]
	return section, nil
}
