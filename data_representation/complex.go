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

// ComplexParams contains the parameters needed for unpacking complex data, including spatially processed complex data.
type ComplexParams struct {
	TotalPoints              int
	DataPoints               int
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
	PrimaryMissingValue      float64
	SecondaryMissingValue    float64
	Bitmap                   *BitmapReader
}

// UnpackComplex unpacks the provided byte slice into a slice of float64 values.
//
// Note: the wgrib2 C codebase on GitHub was especially useful in compiling the logic, in addition to the
// GRIB2 reference documents hosted by NCEP and the regulations hosted by ECWMF.
func (p *ComplexParams) UnpackComplex(packedData []byte) ([]float64, error) {
	g, err := newGroupTracker(p, packedData)
	if err != nil {
		return nil, err
	}

	result := make([]float64, 0, p.TotalPoints)

	var pointIdx int
	for ; pointIdx < p.TotalPoints; pointIdx++ {
		if p.Bitmap.IsMissing(pointIdx) {
			result = append(result, math.NaN())
			continue
		}
		result = append(result, g.nextValue())
	}
	if pointIdx != p.TotalPoints {
		return result, fmt.Errorf("error unpacking complex: expected %d points, got %d", p.TotalPoints, pointIdx)
	}
	return result, nil
}

type groupTracker struct {
	params         *ComplexParams
	currentGroup   int
	pointsInGroup  int
	groupWidth     int
	groupLength    int
	groupRef       int
	currentPoint   int
	currentValue   int
	last1          int
	last2          int
	overallMin     int
	firstFound     bool
	secondFound    bool
	groupRefValues []int
	groupWidths    []int
	groupLengths   []int
	dataStream     *BitStream
}

func newGroupTracker(params *ComplexParams, packedData []byte) (*groupTracker, error) {
	tracker := &groupTracker{
		params: params,
	}
	if params.Order > 2 {
		return nil, fmt.Errorf("error unpacking complex: order %d is not supported", params.Order)
	}

	currentPos := 0
	endPos := (params.Order + 1) * params.SpatialOctets
	tracker.processSpatialOctets(packedData[currentPos:endPos])

	currentPos = endPos
	endPos = currentPos + tracker.bytesRequired(params.BitsPerGroup)
	tracker.groupRefValues = tracker.readUints(params.BitsPerGroup, packedData[currentPos:endPos], 0)

	currentPos = endPos
	endPos = currentPos + tracker.bytesRequired(params.BitsPerGroupWidth)
	tracker.groupWidths = tracker.readUints(params.BitsPerGroupWidth, packedData[currentPos:endPos], params.GroupWidthReference)

	currentPos = endPos
	endPos = currentPos + tracker.bytesRequired(params.BitsPerScaledGroupLength)
	tracker.groupLengths = tracker.readGroupLengths(packedData[currentPos:endPos])
	if err := tracker.checkTotalGroupPoints(); err != nil {
		return nil, err
	}

	tracker.dataStream = NewBitStream(packedData[endPos:])
	tracker.pointsInGroup = tracker.groupLengths[0]
	tracker.groupWidth = tracker.groupWidths[0]
	tracker.groupLength = tracker.groupLengths[0]
	tracker.groupRef = tracker.groupRefValues[0]
	return tracker, nil
}

func (g *groupTracker) nextValue() float64 {
	if g.currentPoint >= g.pointsInGroup {
		g.currentGroup++
		g.pointsInGroup = g.groupLengths[g.currentGroup]
		g.groupWidth = g.groupWidths[g.currentGroup]
		g.groupLength = g.groupLengths[g.currentGroup]
		g.groupRef = g.groupRefValues[g.currentGroup]
		g.currentPoint = 0
	}

	nb := g.groupWidth
	gref := g.groupRef
	val := gref

	if g.params.MissingValueManagement == 0 {
		if nb > 0 {
			val += int(g.dataStream.ReadBits(nb))
		}
	} else if g.params.MissingValueManagement == 1 {
		if nb == 0 {
			m1 := (1 << g.params.BitsPerGroup) - 1
			if m1 == gref {
				val = math.MaxInt64
			}
		} else {
			m1 := (1 << nb) - 1
			dataVal := int(g.dataStream.ReadBits(nb))
			if dataVal == m1 {
				val = math.MaxInt64
			} else {
				val += dataVal
			}
		}
	}
	if g.params.MissingValueManagement == 2 {
		if nb == 0 {
			m1 := (1 << g.params.BitsPerGroup) - 1
			m2 := m1 - 1
			if m1 == gref || m2 == gref {
				val = math.MaxInt64
			}
		} else {
			m1 := (1 << nb) - 1
			m2 := m1 - 1
			dataVal := int(g.dataStream.ReadBits(nb))
			if dataVal == m1 || dataVal == m2 {
				val = math.MaxInt64
			} else {
				val += dataVal
			}
		}
	}

	var result float64
	if g.params.Order == 0 {
		if val == math.MaxInt64 {
			result = math.NaN()
		} else {
			result = u.Unpack(g.params.Ref, val, g.params.BinaryScale, g.params.DecimalScale)
		}
	} else if g.params.Order == 1 {
		if !g.firstFound && val != math.MaxInt64 {
			val = g.last1
			g.firstFound = true
		} else if g.firstFound && val != math.MaxInt64 {
			val += g.last1 + g.overallMin
			g.last1 = val
		}
		if val == math.MaxInt64 {
			result = math.NaN()
		} else {
			result = u.Unpack(g.params.Ref, val, g.params.BinaryScale, g.params.DecimalScale)
		}
	} else if g.params.Order == 2 {
		if !g.firstFound && val != math.MaxInt64 {
			val = g.last2
			g.firstFound = true
		} else if !g.secondFound && val != math.MaxInt64 {
			val = g.last1
			g.secondFound = true
		} else if g.firstFound && g.secondFound && val != math.MaxInt64 {
			val += g.overallMin + g.last1 + (g.last1 - g.last2)
			g.last2 = g.last1
			g.last1 = val
		}
	}
	if val == math.MaxInt64 {
		result = math.NaN()
	} else {
		result = u.Unpack(g.params.Ref, val, g.params.BinaryScale, g.params.DecimalScale)
		if g.params.MissingValueManagement == 1 && math.Abs(result-g.params.PrimaryMissingValue) < 1e-10 {
			result = math.NaN()
		} else if g.params.MissingValueManagement == 2 && (math.Abs(result-g.params.PrimaryMissingValue) < 1e-10 || math.Abs(result-g.params.SecondaryMissingValue) < 1e-10) {
			result = math.NaN()
		}
	}
	g.currentPoint++
	return result
}

func (g *groupTracker) processSpatialOctets(data []byte) {
	var diffMin int
	var initials1, initials2 int
	if g.params.Order > 0 {
		start := 0
		end := g.params.SpatialOctets
		initials1 = UintFromBytes(data[start:end])
		if g.params.Order == 2 {
			start = end
			end = start + g.params.SpatialOctets
			initials2 = UintFromBytes(data[start:end])
		}
		start = end
		end = start + g.params.SpatialOctets
		diffMin = IntFromBytes(data[start:end])
	}
	if g.params.Order == 1 {
		g.last1 = initials1
	} else if g.params.Order == 2 {
		g.last2 = initials1
		g.last1 = initials2
	}
	g.overallMin = diffMin
}

func (g *groupTracker) bytesRequired(bitsRequired int) int {
	return (bitsRequired*g.params.NG + 7) / 8
}

func (g *groupTracker) readUints(bitsRequired int, data []byte, ref int) []int {
	stream := NewBitStream(data)
	values := make([]int, 0, g.params.NG)
	for i := 0; i < g.params.NG; i++ {
		values = append(values, int(stream.ReadBits(bitsRequired))+ref)
	}
	return values
}

func (g *groupTracker) readGroupLengths(data []byte) []int {
	stream := NewBitStream(data)
	values := make([]int, 0, g.params.NG)
	for i := 0; i < g.params.NG; i++ {
		value := int(stream.ReadBits(g.params.BitsPerScaledGroupLength))
		value = g.params.GroupLengthReference + value*g.params.GroupLengthIncrement
		values = append(values, value)
	}
	values[g.params.NG-1] = g.params.LastGroupLength
	return values
}

func (g *groupTracker) checkTotalGroupPoints() error {
	sum := 0
	for _, l := range g.groupLengths {
		sum += l
	}
	if sum != g.params.DataPoints {
		return fmt.Errorf("error unpacking values: total group lengths %d != grid size %d", sum, g.params.DataPoints)
	}
	return nil
}
