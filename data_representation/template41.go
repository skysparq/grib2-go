package data_representation

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
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
	p, err := png.Decode(bytes.NewReader(rec.Data.Data))
	if err != nil {
		return nil, fmt.Errorf("error getting values: %w", err)
	}
	var getValue func(x, y int) int
	switch img := p.(type) {
	case *image.Gray:
		getValue = t.gray8Getter(img)
	case *image.Gray16:
		getValue = t.gray16Getter(img)
	case *image.RGBA:
		if t.BitDepth == 24 {
			getValue = t.rgba24BitGetter(img)
		} else {
			getValue = t.rgba32BitGetter(img)
		}
	default:
		return nil, fmt.Errorf("error getting values: unsupported image type: %T", img)
	}
	bmpR, err := NewBitmapReader(rec)
	if err != nil {
		return nil, fmt.Errorf("error getting values: %w", err)
	}
	return t.getValues(bmpR, p, getValue), nil
}

func (t Template41) getValues(bmpR *BitmapReader, p image.Image, getValue func(x, y int) int) []float64 {
	width, height := p.Bounds().Dx(), p.Bounds().Dy()
	pixels := make([]float64, 0, width*height)
	index := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if bmpR.IsMissing(index) {
				pixels = append(pixels, math.NaN())
				index++
				continue
			}
			value := u.Unpack(t.ReferenceValue, getValue(x, y), t.BinaryScaleFactor, t.DecimalScaleFactor)
			pixels = append(pixels, value)
			index++
		}
	}
	return pixels
}

func (t Template41) rgba24BitGetter(img *image.RGBA) func(x, y int) int {
	return func(x, y int) int {
		r, g, b, _ := img.At(x, y).RGBA()
		r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
		return int(uint32(r8)<<16 | uint32(g8)<<8 | uint32(b8))
	}
}

func (t Template41) rgba32BitGetter(img *image.RGBA) func(x, y int) int {
	return func(x, y int) int {
		r, g, b, a := img.At(x, y).RGBA()
		r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
		return int(uint32(r8)<<24 | uint32(g8)<<16 | uint32(b8)<<8 | uint32(a8))
	}
}

func (t Template41) gray8Getter(img *image.Gray) func(x, y int) int {
	return func(x, y int) int {
		return int(img.GrayAt(x, y).Y)
	}
}

func (t Template41) gray16Getter(img *image.Gray16) func(x, y int) int {
	return func(x, y int) int {
		return int(img.Gray16At(x, y).Y)
	}
}
