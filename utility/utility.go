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
		return -value
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
		return -value
	}
	return value
}

func SignAndMagnitudeInt8(data byte) int {
	negative := (data & 0x80) == 0x80
	data &= 0x7F
	value := int(data)
	if negative {
		return -value
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

func StdLatLngToFloat(value int) float64 {
	return float64(value) * 1e-6
}

func ShiftLongitude(value int) int {
	if value > 180000000 {
		value -= 360000000
	}
	return value
}

func Unpack(ref float64, value int, binaryScale int, decimalScale int) float64 {
	return (ref + (float64(value) * math.Pow(2, float64(binaryScale)))) / math.Pow(10, float64(decimalScale))
}

func UnpackFloat(ref float64, value float64, binaryScale int, decimalScale int) float64 {
	return (ref + (value * math.Pow(2, float64(binaryScale)))) / math.Pow(10, float64(decimalScale))
}
