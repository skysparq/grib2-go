package utility

import (
	"encoding/binary"
	"math"
)

// Int32 returns an int from a big endian 4-byte slice
func Int32(data []byte) int {
	return int(int32(binary.BigEndian.Uint32(data)))
}

// SignAndMagnitudeInt32 returns an int from a big endian 4-byte slice encoded as a sign-and-magnitude integer.
// Unlike two's complement, sign-and-magnitude integers use the first bit to indicate the sign.
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

// SignAndMagnitudeInt16 returns an int from a big endian 2-byte slice encoded as a sign-and-magnitude integer.
// Unlike two's complement, sign-and-magnitude integers use the first bit to indicate the sign.
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

// SignAndMagnitudeInt8 returns an int from a big endian single byte encoded as a sign-and-magnitude integer.
// Unlike two's complement, sign-and-magnitude integers use the first bit to indicate the sign.
func SignAndMagnitudeInt8(data byte) int {
	negative := (data & 0x80) == 0x80
	data &= 0x7F
	value := int(data)
	if negative {
		return -value
	}
	return value
}

// Int64 returns an int from a big endian 8-byte slice
func Int64(data []byte) int {
	return int(binary.BigEndian.Uint64(data))
}

// Uint16 returns an unsigned int from a big endian 2-byte slice
func Uint16(data []byte) int {
	return int(binary.BigEndian.Uint16(data))
}

// Uint32 returns an unsigned int from a big endian 4-byte slice
func Uint32(data []byte) int {
	return int(binary.BigEndian.Uint32(data))
}

// Float32 returns a float32 from a 4-byte slice encoded as a float
func Float32(data []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(data))
}

// StdLatLngToFloat converts a standard latitude/longitude value to a float.
// In grib2 files, latitude and longitudes are stored as 32-bit signed integers with the decimal shifted
// 6 places to the right.
func StdLatLngToFloat(value int) float64 {
	return float64(value) * 1e-6
}

// ShiftLongitude shifts a longitude value to the range [-180, 180].
// By default, they are stored in the range [0, 360]
func ShiftLongitude(value int) int {
	if value > 180000000 {
		value -= 360000000
	}
	return value
}

// Unpack converts a packed value to the original floating point value.
// The formula for recovering packed data is:
//
// Y = (R + X * 2^E) / 10^D
//
// For complex packing:
//
// E = Binary scale factor
//
// D = Decimal scale factor
//
// R = Reference value of the whole field
//
// X = Packed value
func Unpack(ref float64, value int, binaryScale int, decimalScale int) float64 {
	return (ref + (float64(value) * math.Pow(2, float64(binaryScale)))) / math.Pow(10, float64(decimalScale))
}

// UnpackFloat converts a packed floating point value to the original unpacked floating point value.
// The formula for recovering packed data is:
//
// Y = (R + X * 2^E) / 10^D
//
// For complex packing:
//
// E = Binary scale factor
//
// D = Decimal scale factor
//
// R = Reference value of the whole field
//
// X = Packed value
func UnpackFloat(ref float64, value float64, binaryScale int, decimalScale int) float64 {
	return (ref + (value * math.Pow(2, float64(binaryScale)))) / math.Pow(10, float64(decimalScale))
}

func ScaleInt(value int, scale int) float64 {
	return float64(value) / math.Pow(10, float64(scale))
}

// Missing is the value used in code tables to indicate a missing value.
const Missing = 255
