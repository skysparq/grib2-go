package record

import "fmt"

// Section7 contains the data block from section 7 of a GRIB record.
type Section7 struct {
	Length int
	Data   []byte
}

// ParseSection7 parses section 7 of a GRIB record.
func ParseSection7(data SectionData) (section Section7, err error) {
	section.Length = data.Length
	if data.SectionNumber != 7 {
		return section, fmt.Errorf(`error parsing section 7: expected section number 7, got %d`, data.SectionNumber)
	}
	section.Data = data.Bytes[5:]
	return section, nil
}
