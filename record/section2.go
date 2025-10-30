package record

import (
	"fmt"
)

// Section2 contains the fields from section 2 of a GRIB record.
// Section 2 is not guaranteed to be present. If present, it is only presented as a slice of raw bytes straight from the GRIB record.
// No attempt is made to parse the Local Use section. That is left up to the user.
type Section2 struct {
	Length   int
	LocalUse []byte
}

// ParseSection2 parses section 2 of a GRIB record.
func ParseSection2(data SectionData) (section Section2, err error) {
	section.Length = data.Length

	if data.SectionNumber != 2 {
		return section, fmt.Errorf(`error parsing section 2: expected section number 2, got %d`, data.SectionNumber)
	}
	section.LocalUse = data.Bytes[5:]
	return section, nil
}
