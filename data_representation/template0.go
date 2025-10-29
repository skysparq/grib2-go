package data_representation

import (
	"fmt"
	"math"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

type Template0 struct {
	ReferenceValue     float64
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
	t.ReferenceValue = float64(u.Float32(data[0:4]))
	t.BinaryScaleFactor = u.SignAndMagnitudeInt16(data[4:6])
	t.DecimalScaleFactor = u.SignAndMagnitudeInt16(data[6:8])
	t.BitsPerValue = int(data[8])
	t.OriginalFieldType = int(data[9])
	return t, nil
}

func (t Template0) DecimalScale() int {
	return t.DecimalScaleFactor
}

// GetValues in this template uses simple unpacking to retrieve values from the record
func (t Template0) GetValues(rec record.Record) ([]float64, error) {
	getValues := t.getValueReader()
	return getValues(rec, rec.Grid.TotalPoints)
}

func (t Template0) getValueReader() func(rec record.Record, totalPoints int) ([]float64, error) {
	if t.BitsPerValue == 0 {
		return t.unpackConst
	}
	return t.unpackSimple
}

func (t Template0) unpackConst(rec record.Record, totalPoints int) ([]float64, error) {
	ref := u.UnpackFloat(t.ReferenceValue, 0, t.BinaryScaleFactor, t.DecimalScaleFactor)
	values := make([]float64, totalPoints)
	bmpR, err := NewBitmapReader(rec)
	if err != nil {
		return nil, fmt.Errorf("error unpacking const values: %w", err)
	}
	for i := range values {
		if bmpR.IsMissing(i) {
			values[i] = math.NaN()
		} else {
			values[i] = ref
		}
	}
	return values, nil
}

func (t Template0) unpackSimple(rec record.Record, totalPoints int) ([]float64, error) {
	values := make([]float64, totalPoints)
	stream := NewBitStream(rec.Data.Data)
	bmpR, err := NewBitmapReader(rec)
	if err != nil {
		return nil, fmt.Errorf("error unpacking simple values: %w", err)
	}

	for i := range values {
		if bmpR.IsMissing(i) {
			values[i] = math.NaN()
			continue
		}

		packed := int(stream.ReadBits(t.BitsPerValue))
		value := u.Unpack(t.ReferenceValue, packed, t.BinaryScaleFactor, t.DecimalScaleFactor)
		values[i] = value
	}
	stream.Pos()
	return values, nil
}
