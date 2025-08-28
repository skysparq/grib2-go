package utility

import (
	"encoding/binary"
	"math"
)

func Int32(data []byte) int {
	return int(int32(binary.BigEndian.Uint32(data)))
}

func OverflowInt32(data []byte) int {
	value := int(binary.BigEndian.Uint32(data))
	if value > math.MaxInt32 {
		return math.MaxInt32 - value + 1
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
