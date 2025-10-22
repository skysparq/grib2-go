package data_representation_test

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/skysparq/grib2-go/data_representation"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestTemplate41(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordDataDef41)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := data_representation.Template41{}.Parse(rec.DataRepresentation)
	if err != nil {
		t.Fatal(err)
	}
	expected := data_representation.Template41{
		ReferenceValue:     -9990,
		BinaryScaleFactor:  0,
		DecimalScaleFactor: 1,
		BitDepth:           16,
		OriginalFieldType:  0,
	}
	if typed := template.(data_representation.Template41); !reflect.DeepEqual(expected, typed) {
		t.Fatalf("expected\n%+v\nbut got\n%+v", expected, typed)
	}
	encoded, err := json.Marshal(template)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(`%s`, encoded)
}

func TestGetValuesTemplate41(t *testing.T) {
	r, err := os.Open(`../.test_files/MRMS_MergedAzShear_0-2kmAGL_00.50_20251005-112817.grb2`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	rec, err := record.ParseRecord(r, templates.Version33())
	if err != nil {
		t.Fatal(err)
	}
	template, err := data_representation.Template41{}.Parse(rec.DataRepresentation)
	if err != nil {
		t.Fatal(err)
	}
	template41 := template.(data_representation.Template41)
	bitLevel := template.(data_representation.Template41).BitDepth
	t.Logf(`bitLevel: %v`, bitLevel)

	p, err := png.Decode(bytes.NewReader(rec.Data.Data))
	if err != nil {
		t.Fatal(err)
	}
	width, height := p.Bounds().Dx(), p.Bounds().Dy()
	pixels := make([]float32, 0, width*height)
	pmin, pmax := float32(0), float32(0)
	switch img := p.(type) {
	case *image.Gray:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				value := float32(img.GrayAt(x, y).Y)

				ref := getDecimalScaledRef(template41.DecimalScaleFactor, template41.ReferenceValue)
				scale := getScale(template41.DecimalScaleFactor, template41.BinaryScaleFactor)
				value = value * float32(scale)
				value += float32(ref)
				if value < pmin {
					pmin = value
				}
				if value > pmax {
					pmax = value
				}
				pixels = append(pixels, value)
			}
		}
	case *image.Gray16:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				value := float32(img.Gray16At(x, y).Y)

				ref := getDecimalScaledRef(template41.DecimalScaleFactor, template41.ReferenceValue)
				scale := getScale(template41.DecimalScaleFactor, template41.BinaryScaleFactor)
				value = value * float32(scale)
				value += float32(ref)
				if value < pmin {
					pmin = value
				}
				if value > pmax {
					pmax = value
				}
				pixels = append(pixels, value)
			}
		}
	case *image.RGBA:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				var pixel uint32
				if bitLevel == 24 {
					r8, g8, b8 := uint8(r>>8), uint8(g>>8), uint8(b>>8)
					pixel = uint32(r8)<<16 | uint32(g8)<<8 | uint32(b8)
				}
				if bitLevel == 32 {
					r8, g8, b8, a8 := uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8)
					pixel = uint32(r8)<<24 | uint32(g8)<<16 | uint32(b8)<<8 | uint32(a8)
				}

				value := float32(pixel)
				ref := getDecimalScaledRef(template41.DecimalScaleFactor, template41.ReferenceValue)
				scale := getScale(template41.DecimalScaleFactor, template41.BinaryScaleFactor)
				value = value * float32(scale)
				value += float32(ref)
				if value < pmin {
					pmin = value
				}
				if value > pmax {
					pmax = value
				}
				pixels = append(pixels, math.Float32frombits(pixel))
			}
		}
	}

	//e, err := json.Marshal(pixels)
	if err != nil {
		t.Fatal(err)
	}
	//os.WriteFile(`../.test_files/data.json`, e, os.ModePerm)
	t.Logf(`min: %v, max: %v`, pmin, pmax)
}

func getDecimalScaledRef(decimalScaleFactor int, ref float32) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * float64(ref)
}

func getScale(decimalScaleFactor int, binaryScaleFactor int) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * math.Pow(2, float64(binaryScaleFactor))
}
