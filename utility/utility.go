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

func GetDecimalScaledRef(decimalScaleFactor int, ref float32) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * float64(ref)
}

func GetScale(decimalScaleFactor int, binaryScaleFactor int) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * math.Pow(2, float64(binaryScaleFactor))
}

func GetDecimalScaledValue(decimalScaleFactor int, value float64) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * value
}

func NormalizeStdLongitude(value int) float64 {
	if value > 180000000 {
		value -= 360000000
	}
	return float64(value) * 1e-6
}

func NormalizeStdLatitude(value int) float64 {
	return float64(value) * 1e-6
}
