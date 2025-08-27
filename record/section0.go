package record

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Section0 struct {
	Discipline int
	Edition    int
	GribLength int
}

func ParseSection0(r io.Reader) (section0 Section0, err error) {
	sectionBytes, err := readFixedLengthSection(r, 16)
	if err != nil {
		return section0, fmt.Errorf(`error reading section 0: %w`, err)
	}
	if header := string(sectionBytes[0:4]); header != `GRIB` {
		return section0, fmt.Errorf(`error reading section 0: expected header 'GRIB' but got %s`, header)
	}
	section0.Discipline = int(sectionBytes[6])
	section0.Edition = int(sectionBytes[7])
	section0.GribLength = int(binary.BigEndian.Uint64(sectionBytes[8:16]))
	return section0, nil
}
