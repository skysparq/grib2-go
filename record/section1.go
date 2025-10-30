package record

import (
	"fmt"
	"time"

	u "github.com/skysparq/grib2-go/utility"
)

// Section1 contains the fields from section 1 of a GRIB record.
type Section1 struct {
	Length                    int
	OriginatingCenter         int
	OriginatingSubCenter      int
	MasterTableVersion        int
	LocalTableVersion         int
	ReferenceTimeSignificance int
	Year                      int
	Month                     int
	Day                       int
	Hour                      int
	Minute                    int
	Second                    int
	ProductionStatus          int
	DataType                  int
	Reserved                  []byte
}

// ParseSection1 parses section 1 of a GRIB record.
func ParseSection1(data SectionData) (section Section1, err error) {
	section.Length = data.Length
	if data.SectionNumber != 1 {
		return section, fmt.Errorf(`error parsing section 1: expected section number 1, got %d`, data.SectionNumber)
	}
	section.OriginatingCenter = u.Uint16(data.Bytes[5:7])
	section.OriginatingSubCenter = u.Uint16(data.Bytes[7:9])
	section.MasterTableVersion = int(data.Bytes[9])
	section.LocalTableVersion = int(data.Bytes[10])
	section.ReferenceTimeSignificance = int(data.Bytes[11])
	section.Year = u.Uint16(data.Bytes[12:14])
	section.Month = int(data.Bytes[14])
	section.Day = int(data.Bytes[15])
	section.Hour = int(data.Bytes[16])
	section.Minute = int(data.Bytes[17])
	section.Second = int(data.Bytes[18])
	section.ProductionStatus = int(data.Bytes[19])
	section.DataType = int(data.Bytes[20])
	section.Reserved = data.Bytes[21:]

	return section, nil
}

// Time returns the cycle time of the GRIB record.
// For observational data, this is typically the time the record was generated
// For forecast data, this is typically the starting time of the forecast cycle.
func (s Section1) Time() time.Time {
	return time.Date(s.Year, time.Month(s.Month), s.Day, s.Hour, s.Minute, s.Second, 0, time.UTC)
}
