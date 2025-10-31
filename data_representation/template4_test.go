package data_representation_test

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
)

func TestUnpackTemplate4with32bits(t *testing.T) {
	dataDef, err := data_representation.Template4{}.Parse(
		record.Section5{
			DataRepresentationTemplateNumber: 4,
			DataRepresentationTemplateData:   []byte{1},
			Templates:                        templates.Version33(),
		})
	if err != nil {
		t.Fatal(err)
	}

	blob := make([]byte, 0, 12)
	blob = binary.BigEndian.AppendUint32(blob, math.Float32bits(1.0))
	blob = binary.BigEndian.AppendUint32(blob, math.Float32bits(2.0))
	blob = binary.BigEndian.AppendUint32(blob, math.Float32bits(3.0))

	values, err := dataDef.GetValues(record.Record{
		Grid:   record.Section3{TotalPoints: 3},
		BitMap: record.Section6{BitmapIndicator: 255},
		Data:   record.Section7{Data: blob},
	})
	if err != nil {
		t.Fatal(err)
	}
	if expected := 3; len(values) != expected {
		t.Fatalf(`expected %v values but got %v`, expected, len(values))
	}
	if values[0] != 1.0 {
		t.Fatalf(`expected %v but got %v`, 1.0, values[0])
	}
	if values[1] != 2.0 {
		t.Fatalf(`expected %v but got %v`, 2.0, values[1])
	}
	if values[2] != 3.0 {
		t.Fatalf(`expected %v but got %v`, 3.0, values[2])
	}
}

func TestUnpackTemplate4with64bits(t *testing.T) {
	dataDef, err := data_representation.Template4{}.Parse(
		record.Section5{
			DataRepresentationTemplateNumber: 4,
			DataRepresentationTemplateData:   []byte{2},
			Templates:                        templates.Version33(),
		})
	if err != nil {
		t.Fatal(err)
	}

	blob := make([]byte, 0, 24)
	blob = binary.BigEndian.AppendUint64(blob, math.Float64bits(1.0))
	blob = binary.BigEndian.AppendUint64(blob, math.Float64bits(2.0))
	blob = binary.BigEndian.AppendUint64(blob, math.Float64bits(3.0))

	values, err := dataDef.GetValues(record.Record{
		Grid:   record.Section3{TotalPoints: 3},
		BitMap: record.Section6{BitmapIndicator: 255},
		Data:   record.Section7{Data: blob},
	})
	if err != nil {
		t.Fatal(err)
	}
	if expected := 3; len(values) != expected {
		t.Fatalf(`expected %v values but got %v`, expected, len(values))
	}
	if values[0] != 1.0 {
		t.Fatalf(`expected %v but got %v`, 1.0, values[0])
	}
	if values[1] != 2.0 {
		t.Fatalf(`expected %v but got %v`, 2.0, values[1])
	}
	if values[2] != 3.0 {
		t.Fatalf(`expected %v but got %v`, 3.0, values[2])
	}
}
