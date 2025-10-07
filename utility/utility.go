package utility

import (
	"encoding/binary"
	"math"
)

func Int32(data []byte) int {
	return int(int32(binary.BigEndian.Uint32(data)))
}

func SignAndMagnitudeInt32(data []byte) int {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	negative := (dataCopy[0] & 0x80) == 0x80
	dataCopy[0] &= 0x7F
	value := int(binary.BigEndian.Uint32(dataCopy))
	if negative {
		value *= -1
	}
	return value
}

func SignAndMagnitudeInt16(data []byte) int {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	negative := (dataCopy[0] & 0x80) == 0x80
	dataCopy[0] &= 0x7F
	value := int(binary.BigEndian.Uint16(dataCopy))
	if negative {
		value *= -1
	}
	return value
}

func Int64(data []byte) int {
	return int(binary.BigEndian.Uint64(data))
}

func Uint16(data []byte) int {
	return int(binary.BigEndian.Uint16(data))
}

func Uint32(data []byte) int {
	return int(binary.BigEndian.Uint32(data))
}

func Float32(data []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(data))
}
