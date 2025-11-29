package data_representation

import (
	"encoding/binary"
	"fmt"
	"iter"
	"math"

	"github.com/skysparq/grib2-go/record"
)

// Template4 contains the fields for Grid Point Data - IEEE Floating Point Data
type Template4 struct {
	Precision int
}

// Parse fills in the template from the provided section
func (t Template4) Parse(section record.Section5) (record.DataRepresentationDefinition, error) {
	err := checkSectionNum(section, 4)
	if err != nil {
		return t, err
	}

	data := section.DataRepresentationTemplateData
	t.Precision = int(data[0])

	return t, nil
}

// DecimalScale returns the decimal scale factor of the record. The decimal scale factor is used to shift the
// decimal point of a decoded value to the correct position.
func (t Template4) DecimalScale() int {
	return 0
}

// GetValues unpacks the record's data into the original values
func (t Template4) GetValues(rec record.Record) ([]float64, error) {
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

func (t Template4) ValuesIterator(rec record.Record) (iter.Seq2[int, float64], error) {
	var get func([]byte, int, int) float64
	var increment int
	if t.Precision == 1 {
		get = t.readFloat32
		increment = 4
	} else if t.Precision == 2 {
		get = t.readFloat64
		increment = 8
	} else {
		return nil, fmt.Errorf(`error getting values iterator: unsupported precision: %d`, t.Precision)
	}

	bitmap := NewBitmapReader(rec)
	totalPoints := rec.Grid.TotalPoints
	blob := rec.Data.Data

	return func(yield func(int, float64) bool) {
		i := 0
		for {
			if i >= totalPoints {
				return
			}
			var value float64
			if bitmap.IsMissing(i) {
				value = math.NaN()
			} else {
				value = get(blob, i, increment)
			}
			if !yield(i, value) {
				return
			}
			i++
		}
	}, nil
}

func (t Template4) readFloat32(data []byte, index, increment int) float64 {
	start := index * increment
	end := start + increment
	value := float64(math.Float32frombits(binary.BigEndian.Uint32(data[start:end])))
	return value
}

func (t Template4) readFloat64(data []byte, index, increment int) float64 {
	start := index * increment
	end := start + increment
	value := math.Float64frombits(binary.BigEndian.Uint64(data[start:end]))
	return value
}
