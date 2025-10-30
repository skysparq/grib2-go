package record

import (
	"encoding/binary"
	"fmt"
	"io"
)

// SectionData contains the raw bytes from a section of a GRIB record.
// Each section begins by specifying the length of the section and its section number.
// This struct is used as an intermediary step to read each section.
type SectionData struct {
	Length        int
	SectionNumber int
	Bytes         []byte
}

func readVariableLengthSection(r io.Reader) (data SectionData, err error) {
	lengthBytes := make([]byte, 4)
	totalRead, err := io.ReadFull(r, lengthBytes)
	if err != nil {
		return data, err
	}
	if totalRead != 4 {
		return data, fmt.Errorf(`expected 4 bytes read for length, got %d`, totalRead)
	}
	length := int(binary.BigEndian.Uint32(lengthBytes))
	data.Length = length

	remainderLength := length - 4
	remainderBytes := make([]byte, remainderLength)
	totalRead, err = io.ReadFull(r, remainderBytes)
	if err != nil {
		return data, err
	}
	if totalRead != remainderLength {
		return data, fmt.Errorf(`expected %d bytes read for rest of record, got %d`, remainderLength, totalRead)
	}
	data.SectionNumber = int(remainderBytes[0])
	data.Bytes = append(lengthBytes, remainderBytes...)
	return data, nil
}

func readFixedLengthSection(r io.Reader, length int) ([]byte, error) {
	sectionBytes := make([]byte, length)
	totalRead, err := io.ReadFull(r, sectionBytes)
	if err != nil {
		return sectionBytes, err
	}
	if totalRead != length {
		return sectionBytes, fmt.Errorf(`expected %d bytes read, got %d`, length, totalRead)
	}
	return sectionBytes, nil
}
