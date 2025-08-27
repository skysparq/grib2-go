package record

import (
	"fmt"
	"io"
)

type Record struct {
	Indicator         Section0
	Identification    Section1
	LocalUse          Section2
	GridDefinition    Section3
	ProductDefinition Section4
	BitMap            Section5
	Data              Section6
}

func ParseRecord(r io.Reader) (record Record, err error) {
	record.Indicator, err = ParseSection0(r)
	if err != nil {
		return record, fmt.Errorf("error parsing record: %w", err)
	}
	record.Identification, err = ParseSection1(r)
	if err != nil {
		return record, fmt.Errorf("error parsing record: %w", err)
	}
	return record, nil
}
