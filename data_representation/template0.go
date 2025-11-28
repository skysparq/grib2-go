package data_representation

import (
	"iter"
	"math"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

// Template0 contains the fields for Grid Point Data - Simple Packing
type Template0 struct {
	ReferenceValue     float64
	BinaryScaleFactor  int
	DecimalScaleFactor int
	BitsPerValue       int
	OriginalFieldType  int
}

// Parse fills in the template from the provided section
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

// DecimalScale returns the decimal scale factor of the record. The decimal scale factor is used to shift the
// decimal point of a decoded value to the correct position.
func (t Template0) DecimalScale() int {
	return t.DecimalScaleFactor
}

// GetValues unpacks the record's data into the original values
func (t Template0) GetValues(rec record.Record) ([]float64, error) {
	getValues := t.getValueReader()
	return getValues(rec, rec.Grid.TotalPoints), nil
}

func (t Template0) ValuesIterator(rec record.Record) (iter.Seq2[int, float64], error) {
	if t.BitsPerValue == 0 {
		return t.constIterator(rec, rec.Grid.TotalPoints), nil
	}
	return t.simpleIterator(rec, rec.Grid.TotalPoints), nil
}

func (t Template0) getValueReader() func(rec record.Record, totalPoints int) []float64 {
	if t.BitsPerValue == 0 {
		return t.unpackConst
	}
	return t.unpackSimple
}

func (t Template0) unpackConst(rec record.Record, totalPoints int) []float64 {
	ref := u.UnpackFloat(t.ReferenceValue, 0, t.BinaryScaleFactor, t.DecimalScaleFactor)
	values := make([]float64, totalPoints)
	bmpR := NewBitmapReader(rec)
	for i := range values {
		if bmpR.IsMissing(i) {
			values[i] = math.NaN()
		} else {
			values[i] = ref
		}
	}
	return values
}

func (t Template0) unpackSimple(rec record.Record, totalPoints int) []float64 {
	values := make([]float64, totalPoints)
	stream := NewBitStream(rec.Data.Data)
	bmpR := NewBitmapReader(rec)

	for i := range values {
		if bmpR.IsMissing(i) {
			values[i] = math.NaN()
			continue
		}

		packed := int(stream.ReadBits(t.BitsPerValue))
		value := u.Unpack(t.ReferenceValue, packed, t.BinaryScaleFactor, t.DecimalScaleFactor)
		values[i] = value
	}
	return values
}

func (t Template0) simpleIterator(rec record.Record, totalPoints int) iter.Seq2[int, float64] {
	return func(yield func(int, float64) bool) {
		i := 0
		stream := NewBitStream(rec.Data.Data)
		bmpR := NewBitmapReader(rec)
		for {
			if i >= totalPoints {
				return
			}
			var value float64
			if bmpR.IsMissing(i) {
				value = math.NaN()
			} else {
				packed := int(stream.ReadBits(t.BitsPerValue))
				value = u.Unpack(t.ReferenceValue, packed, t.BinaryScaleFactor, t.DecimalScaleFactor)
			}
			if !yield(i, value) {
				return
			}
			i++
		}
	}
}

func (t Template0) constIterator(rec record.Record, totalPoints int) iter.Seq2[int, float64] {
	ref := u.UnpackFloat(t.ReferenceValue, 0, t.BinaryScaleFactor, t.DecimalScaleFactor)

	return func(yield func(int, float64) bool) {
		i := 0
		bmpR := NewBitmapReader(rec)

		for {
			if i >= totalPoints {
				return
			}
			var value float64
			if bmpR.IsMissing(i) {
				value = math.NaN()

			} else {
				value = ref
			}
			if !yield(i, value) {
				return
			}
			i++
		}
	}
}
