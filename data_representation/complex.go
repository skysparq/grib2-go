package data_representation

import (
	"fmt"
	"math"

	u "github.com/skysparq/grib2-go/utility"
)

// The formula for recovering packed data is:
// Y = (R + (X1 + X2) * 2^E) / 10^D
//
// For complex packing:
// E = Binary scale factor
// D = Decimal scale factor
// R = Reference value of the whole field
// X1 = Reference value (scaled integer) of the group the data value belongs to
// X2 = Scaled value with the group reference value removed

type ComplexParams struct {
	TotalPoints              int
	Order                    int
	SpatialOctets            int
	NG                       int
	BitsPerGroup             int
	BitsPerGroupWidth        int
	BitsPerScaledGroupLength int
	GroupWidthReference      int
	GroupLengthReference     int
	GroupLengthIncrement     int
	LastGroupLength          int
	Ref                      float64
	BinaryScale              int
	DecimalScale             int
	MissingValueManagement   int
	Bitmap                   BitmapReader

	initials       []int
	overallMin     int
	groupRefValues []int
	groupWidths    []int
	groupLengths   []int
	packedValues   []byte
}

func (p ComplexParams) UnpackComplex(packedData []byte) ([]float64, error) {
	if p.Order > 2 {
		return nil, fmt.Errorf("error unpacking complex: order %d is not supported", p.Order)
	}

	currentPos := 0
	endPos := (p.Order + 1) * p.SpatialOctets
	p.initials, p.overallMin = p.processSpatialOctets(packedData[currentPos:endPos])

	currentPos = endPos
	endPos = currentPos + p.bytesRequired(p.BitsPerGroup)
	p.groupRefValues = p.readInts(p.BitsPerGroup, packedData[currentPos:endPos], 0)

	currentPos = endPos
	endPos = currentPos + p.bytesRequired(p.BitsPerGroupWidth)
	p.groupWidths = p.readUInts(p.BitsPerGroupWidth, packedData[currentPos:endPos], p.GroupWidthReference)

	currentPos = endPos
	endPos = currentPos + p.bytesRequired(p.BitsPerScaledGroupLength)
	p.groupLengths = p.readGroupLengths(packedData[currentPos:endPos])
	if err := p.checkTotalGroupPoints(p.groupLengths); err != nil {
		return nil, err
	}

	p.packedValues = packedData[endPos:]
	dataStream := NewBitStream(p.packedValues)

	pointIdx := 0
	var last1, last2 int
	if p.Order == 1 {
		last1 = p.initials[0]
	} else if p.Order == 2 {
		last2 = p.initials[0]
		last1 = p.initials[1]
	}
	result := make([]float64, 0, p.TotalPoints)
	firstFound, secondFound := false, false

	for g := 0; g < p.NG; g++ {
		nb := p.groupWidths[g]
		nl := p.groupLengths[g]
		gref := p.groupRefValues[g]

		for i := 0; i < nl; i++ {
			val := gref

			if p.MissingValueManagement == 0 {
				if nb > 0 {
					val += int(dataStream.ReadBits(nb))
				}
			}
			if p.MissingValueManagement == 1 {
				if nb == 0 {
					m1 := (1 << p.BitsPerGroup) - 1
					if m1 == gref {
						val = math.MaxInt64
					}
				} else {
					m1 := (1 << (nb - 1)) - 1
					dataVal := int(dataStream.ReadBits(nb))
					if dataVal == m1 {
						val = math.MaxInt64
					} else {
						val += dataVal
					}
				}
			}
			if p.MissingValueManagement == 2 {
				if nb == 0 {
					m1 := (1 << p.BitsPerGroup) - 1
					m2 := m1 - 1
					if m1 == gref || m2 == gref {
						val = math.MaxInt64
					}
				} else {
					m1 := (1 << (nb - 1)) - 1
					m2 := m1 - 1
					dataVal := int(dataStream.ReadBits(nb))
					if dataVal == m1 || dataVal == m2 {
						val = math.MaxInt64
					} else {
						val += dataVal
					}
				}
			}

			if p.Order == 1 {
				if !firstFound && val != math.MaxInt64 {
					val = p.initials[0]
					firstFound = true
				}
				if val != math.MaxInt64 {
					val += last1 + p.overallMin
					last1 = val
				}
			}
			if p.Order == 2 {
				if !firstFound && val != math.MaxInt64 {
					val = p.initials[0]
					firstFound = true
				} else if !secondFound && val != math.MaxInt64 {
					val = p.initials[1]
					secondFound = true
				}
				if val != math.MaxInt64 {
					val += p.overallMin + last1 + (last1 - last2)
					last2 = last1
					last1 = val
				}
			}

			if p.Bitmap.IsSet(pointIdx) || val == math.MaxInt64 {
				result = append(result, math.NaN())
			} else {
				result = append(result, u.Unpack(p.Ref, val, p.BinaryScale, p.DecimalScale))
			}
			pointIdx++
		}
	}

	if pointIdx != len(result) {
		return result, fmt.Errorf("error unpacking complex: decoded %d points, expected %d", pointIdx, len(result))
	}

	return result, nil
}

func (p ComplexParams) processSpatialOctets(data []byte) ([]int, int) {
	stream := NewBitStream(data)
	var diffMin int
	initials := make([]int, 0, p.Order)
	if p.Order > 0 {
		initials = append(initials, int(stream.ReadSignedBits(p.SpatialOctets*8)))
		if p.Order == 2 {
			initials = append(initials, int(stream.ReadSignedBits(p.SpatialOctets*8)))
		}
		diffMin = int(stream.ReadSignedBits(p.SpatialOctets * 8))
	}
	return initials, diffMin
}

func (p ComplexParams) bytesRequired(bitsRequired int) int {
	return (bitsRequired*p.NG + 7) / 8
}

func (p ComplexParams) readInts(bitsRequired int, data []byte, ref int) []int {
	stream := NewBitStream(data)
	values := make([]int, 0, p.NG)
	for i := 0; i < p.NG; i++ {
		values = append(values, int(stream.ReadSignedBits(bitsRequired))+ref)
	}
	return values
}

func (p ComplexParams) readUInts(bitsRequired int, data []byte, ref int) []int {
	stream := NewBitStream(data)
	values := make([]int, 0, p.NG)
	for i := 0; i < p.NG; i++ {
		values = append(values, int(stream.ReadBits(bitsRequired))+ref)
	}
	return values
}

func (p ComplexParams) readGroupLengths(data []byte) []int {
	stream := NewBitStream(data)
	values := make([]int, 0, p.NG)
	for i := 0; i < p.NG; i++ {
		value := int(stream.ReadBits(p.BitsPerScaledGroupLength))
		value = p.GroupLengthReference + value*p.GroupLengthIncrement
		values = append(values, value)
	}
	values[p.NG-1] = p.LastGroupLength
	return values
}

func (p ComplexParams) checkTotalGroupPoints(groupLengths []int) error {
	sum := 0
	for _, l := range groupLengths {
		sum += l
	}
	if sum != p.TotalPoints {
		return fmt.Errorf("error unpacking values: total group lengths %d != grid size %d", sum, p.TotalPoints)
	}
	return nil
}
