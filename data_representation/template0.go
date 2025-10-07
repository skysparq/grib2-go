package data_representation

import (
	"bytes"
	"fmt"
	"math"

	"github.com/icza/bitio"
	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template0 struct {
	ReferenceValue     float32
	BinaryScaleFactor  int
	DecimalScaleFactor int
	BitsPerValue       int
	OriginalFieldType  int
}

func (t Template0) Parse(section record.Section5) (record.DataRepresentationDefinition, error) {
	err := checkSectionNum(section, 0)
	if err != nil {
		return t, err
	}

	data := section.DataRepresentationTemplateData
	t.ReferenceValue = u.Float32(data[0:4])
	t.BinaryScaleFactor = u.SignAndMagnitudeInt16(data[4:6])
	t.DecimalScaleFactor = u.SignAndMagnitudeInt16(data[6:8])
	t.BitsPerValue = int(data[8])
	t.OriginalFieldType = int(data[9])
	return t, nil
}

// GetValues in this template uses simple unpacking to retrieve values from the record
func (t Template0) GetValues(rec record.Record) ([]float32, error) {
	getValues := t.getValueReader()
	return getValues(rec.Data.Data, rec.GridDefinition.TotalPoints)
}

func (t Template0) getValueReader() func(data []byte, totalPoints int) ([]float32, error) {
	if t.BitsPerValue == 0 {
		return t.unpackConst
	}
	return t.unpackSimple
}

func (t Template0) unpackConst(_ []byte, totalPoints int) ([]float32, error) {
	ref := getDecimalScaledRef(t.DecimalScaleFactor, t.ReferenceValue)
	values := make([]float32, totalPoints)
	for i := range values {
		values[i] = float32(ref)
	}
	return values, nil
}

func (t Template0) unpackSimple(data []byte, totalPoints int) ([]float32, error) {
	values := make([]float32, totalPoints)
	ref := getDecimalScaledRef(t.DecimalScaleFactor, t.ReferenceValue)
	scale := getScale(t.DecimalScaleFactor, t.BinaryScaleFactor)
	reader := bitio.NewReader(bytes.NewBuffer(data))
	for i := range values {
		packed, err := reader.ReadBits(uint8(t.BitsPerValue))
		if err != nil {
			return nil, fmt.Errorf(`error performing simple unpack with bitmap for value %d: %w`, i, err)
		}

		value := float64(math.Float32frombits(uint32(packed)))
		value = value * scale
		value += ref
		values[i] = float32(value)
	}
	return values, nil
}
