package record

import (
	"encoding/binary"
	"fmt"
	"io"
)

type SectionData struct {
	Length int
	Data   []byte
}

func readVariableLengthSection(r io.Reader) (data SectionData, err error) {
	lengthBytes := make([]byte, 4)
	totalRead, err := r.Read(lengthBytes)
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
	totalRead, err = r.Read(remainderBytes)
	if err != nil {
		return data, err
	}
	if totalRead != remainderLength {
		return data, fmt.Errorf(`expected %d bytes read for rest of record, got %d`, remainderLength, totalRead)
	}
	return data, nil
}

func readFixedLengthSection(r io.Reader, length int) ([]byte, error) {
	sectionBytes := make([]byte, length)
	totalRead, err := r.Read(sectionBytes)
	if err != nil {
		return sectionBytes, err
	}
	if totalRead != length {
		return sectionBytes, fmt.Errorf(`expected %d bytes read, got %d`, length, totalRead)
	}
	return sectionBytes, nil
}
