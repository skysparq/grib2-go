package data_representation

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"iter"
	"math"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template41 contains the fields for Grid point data - Portable Network Graphics (PNG)
type Template41 struct {
	ReferenceValue     float64
	BinaryScaleFactor  int
	DecimalScaleFactor int
	BitDepth           int
	OriginalFieldType  int
}

// Parse fills in the template from the provided section
func (t Template41) Parse(section record.Section5) (record.DataRepresentationDefinition, error) {
	err := checkSectionNum(section, 41)
	if err != nil {
		return t, err
	}

	data := section.DataRepresentationTemplateData
	t.ReferenceValue = float64(u.Float32(data[0:4]))
	t.BinaryScaleFactor = u.SignAndMagnitudeInt16(data[4:6])
	t.DecimalScaleFactor = u.SignAndMagnitudeInt16(data[6:8])
	t.BitDepth = int(data[8])
	t.OriginalFieldType = int(data[9])
	return t, nil
}

// DecimalScale returns the decimal scale factor of the record. The decimal scale factor is used to shift the
// decimal point of a decoded value to the correct position.
func (t Template41) DecimalScale() int {
	return t.DecimalScaleFactor
}

// GetValues unpacks the record's data into the original values
func (t Template41) GetValues(rec record.Record) ([]float64, error) {
	iterator, err := t.ValuesIterator(rec)
	if err != nil {
		return nil, fmt.Errorf("error getting values: %w", err)
	}
	values := make([]float64, rec.Grid.TotalPoints)
	for i, v := range iterator {
		values[i] = v
	}
	return values, nil
}

func (t Template41) ValuesIterator(rec record.Record) (iter.Seq2[int, float64], error) {
	p, err := png.Decode(bytes.NewReader(rec.Data.Data))
	if err != nil {
		return nil, fmt.Errorf("error getting values iterator: %w", err)
	}

	var getValue func(i int) int
	switch img := p.(type) {
	case *image.Gray:
		getValue = t.gray8Getter(img.Pix)
	case *image.Gray16:
		getValue = t.gray16Getter(img.Pix)
	case *image.RGBA:
		if t.BitDepth == 24 {
			getValue = t.rgba24BitGetter(img.Pix)
		} else {
			getValue = t.rgba32BitGetter(img.Pix)
		}
	default:
		return nil, fmt.Errorf("error getting values iterator: unsupported image type: %T", img)
	}
	bmpR := NewBitmapReader(rec)
	if err != nil {
		return nil, fmt.Errorf("error getting values iterator: %w", err)
	}
	return func(yield func(int, float64) bool) {
		i := 0
		for {
			if i >= rec.Grid.TotalPoints {
				return
			}

			var value float64
			if bmpR.IsMissing(i) {
				value = math.NaN()
			} else {
				value = u.Unpack(t.ReferenceValue, getValue(i), t.BinaryScaleFactor, t.DecimalScaleFactor)
			}
			if !yield(i, value) {
				return
			}
			i++
		}
	}, nil
}

func (t Template41) rgba24BitGetter(bytes []uint8) func(i int) int {
	return func(i int) int {
		return int(uint32(bytes[i*3])<<16 | uint32(bytes[i*3+1])<<8 | uint32(bytes[i*3+2]))
	}
}

func (t Template41) rgba32BitGetter(bytes []uint8) func(i int) int {
	return func(i int) int {
		return int(binary.LittleEndian.Uint32(bytes[i*4 : i*4+4]))
	}
}

func (t Template41) gray8Getter(bytes []uint8) func(i int) int {
	return func(i int) int {
		return int(bytes[i])
	}
}

func (t Template41) gray16Getter(bytes []uint8) func(i int) int {
	return func(i int) int {
		return int(binary.LittleEndian.Uint16(bytes[i*2 : i*2+2]))
	}
}
