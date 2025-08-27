package record

import (
	"fmt"
	"io"
)

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

func ParseSection1(r io.Reader) (section1 Section1, err error) {
	data, err := readVariableLengthSection(r)
	if err != nil {
		return section1, fmt.Errorf(`error parsing section 1: %w`, err)
	}
	section1.Length = data.Length

	return section1, nil
}
