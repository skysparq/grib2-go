package record_test

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

var template = templates.Version33()

func TestParseGfsRecordPointInTime(t *testing.T) {
	_, r, err := test_files.Load(test_files.SingleRecordProdDef0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r, template)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection0(rec.Indicator)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection1(rec.Identification)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection2(rec.LocalUse)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection3(rec.Grid)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection4(rec.Product)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection5(rec.DataRepresentation)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection6(rec.BitMap)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection7(rec.Data)
	if err != nil {
		t.Fatal(err)
	}
}

func checkSection0(section record.Section0) error {
	if expected := 0; section.Discipline != 0 {
		return fmt.Errorf(`expected discipline %d, got %d`, expected, section.Discipline)
	}
	if expected := 2; section.Edition != 2 {
		return fmt.Errorf(`expected edition %d, got %d`, expected, section.Edition)
	}
	if expected := 1_000_914; section.GribLength != expected {
		return fmt.Errorf(`expected grib length %d, got %d`, expected, section.GribLength)
	}
	return nil
}

func checkSection1(section record.Section1) error {
	expected := record.Section1{
		Length:                    21,
		OriginatingCenter:         7,
		OriginatingSubCenter:      0,
		MasterTableVersion:        2,
		LocalTableVersion:         1,
		ReferenceTimeSignificance: 1,
		Year:                      2025,
		Month:                     3,
		Day:                       5,
		Hour:                      6,
		Minute:                    0,
		Second:                    0,
		ProductionStatus:          0,
		DataType:                  1,
		Reserved:                  []byte{},
	}
	if !reflect.DeepEqual(expected, section) {
		return fmt.Errorf("expected\n%+v\nbut got\n%+v", expected, section)
	}
	if date := time.Date(2025, 3, 5, 6, 0, 0, 0, time.UTC); date != section.Time() {
		return fmt.Errorf(`expected time %v, got %v`, date, section.Time())
	}
	return nil
}

func checkSection2(section record.Section2) error {
	if expected := 0; section.Length != expected {
		return fmt.Errorf(`expected length %d, got %d`, expected, section.Length)
	}
	return nil
}

func checkSection3(section record.Section3) error {
	expected := record.Section3{
		Length:                       72,
		GridSourceDefinition:         0,
		TotalPoints:                  1_038_240,
		OctetsForOptionalPointList:   0,
		InterpretationOfPointList:    0,
		GridDefinitionTemplateNumber: 0,
		GridDefinitionTemplateData:   section.GridDefinitionTemplateData,
		OptionalPointListData:        []byte{},
		Templates:                    template,
	}
	if !reflect.DeepEqual(expected, section) {
		return fmt.Errorf("expected\n%+v\nbut got\n%+v", expected, section)
	}
	if expectedLen := 72 - 15 + 1; len(section.GridDefinitionTemplateData) != expectedLen {
		return fmt.Errorf(`expected length %v, got %d`, expectedLen, len(section.GridDefinitionTemplateData))
	}
	return nil
}

func checkSection4(section record.Section4) error {
	expected := record.Section4{
		Length:                          34,
		CoordinateValuesAfterTemplate:   0,
		ProductDefinitionTemplateNumber: 0,
		ProductDefinitionTemplateData:   section.ProductDefinitionTemplateData,
		CoordinateValuesData:            []byte{},
		Templates:                       template,
	}
	if !reflect.DeepEqual(expected, section) {
		return fmt.Errorf("expected\n%+v\nbut got\n%+v", expected, section)
	}
	if expectedLen := 25; len(section.ProductDefinitionTemplateData) != expectedLen {
		return fmt.Errorf(`expected length %v, got %d`, expectedLen, len(section.ProductDefinitionTemplateData))
	}
	return nil
}

func checkSection5(section record.Section5) error {
	expected := record.Section5{
		Length:                           49,
		TotalDataPoints:                  1_038_240,
		DataRepresentationTemplateNumber: 3,
		DataRepresentationTemplateData:   section.DataRepresentationTemplateData,
		Templates:                        template,
	}
	if !reflect.DeepEqual(expected, section) {
		return fmt.Errorf("expected\n%+v\nbut got\n%+v", expected, section)
	}
	if expectedLen := 38; len(section.DataRepresentationTemplateData) != expectedLen {
		return fmt.Errorf(`expected data representation length %d, got %d`, expectedLen, len(section.DataRepresentationTemplateData))
	}
	return nil
}

func checkSection6(section record.Section6) error {
	if expected := 6; section.Length != expected {
		return fmt.Errorf(`expected length %d, got %d`, expected, section.Length)
	}
	if expected := 255; section.BitmapIndicator != expected {
		return fmt.Errorf(`expected bitmap indicator %d, got %d`, expected, section.BitmapIndicator)
	}
	if expected := 0; len(section.BitmapData) != expected {
		return fmt.Errorf(`expected data representation length %d, got %d`, expected, len(section.BitmapData))
	}
	return nil
}

func checkSection7(section record.Section7) error {
	if expected := 1_000_712; section.Length != expected {
		return fmt.Errorf(`expected length %d, got %d`, expected, section.Length)
	}
	if expected := 1_000_707; len(section.Data) != expected {
		return fmt.Errorf(`expected data representation length %d, got %d`, expected, len(section.Data))
	}
	return nil
}

func TestParseGfsRecordAccumulatedOverTime(t *testing.T) {
	v33 := templates.Version33()
	_, r, err := test_files.Load(test_files.SingleRecordProdDef0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	_, err = record.ParseRecord(r, v33)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGfsNoErrors(t *testing.T) {
	v33 := templates.Version33()
	_, r, err := test_files.Load(test_files.FullGfsFile)
	if err != nil {
		t.Fatal(err)
	}

	f := file.NewGribFile(r, v33)
	err = testFile(f)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHrrrNoErrors(t *testing.T) {
	v33 := templates.Version33()
	_, r, err := test_files.Load(`.test_files/hrrr.t00z.wrfnatf01.grib2`)
	if err != nil {
		t.Fatal(err)
	}

	f := file.NewGribFile(r, v33)
	err = testFile(f)
	if err != nil {
		t.Fatal(err)
	}
}

func testFile(f file.GribFile) error {
	recNum := 0
	for rec, err := range f.Records {
		recNum++
		if err != nil {
			return fmt.Errorf(`error retrieving record for record %v: %w`, recNum, err)
		}
		def, err := rec.DataRepresentation.Definition()
		if err != nil {
			return fmt.Errorf(`error retrieving data representation definition for record %v: %w`, recNum, err)
		}
		_, err = def.GetValues(rec)
		if err != nil {
			return fmt.Errorf(`error getting values for record %v: %w`, recNum, err)
		}
	}
	return nil
}

type floatReader struct {
	r      io.Reader
	buffer [4]byte
}

func (f *floatReader) Next() (float32, error) {
	_, err := io.ReadFull(f.r, f.buffer[:])
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(f.buffer[:])), nil
}

func TestGfsValuesToEccodes(t *testing.T) {
	t.Skip(`This is a long-running validation against values generated from a full GFS grib file using eccodes. It will take several minutes to run. The expected values can be downloaded from https://drive.google.com/file/d/1MhQ1EVHNZsaLBZZYO1ziUDpOfy3t7viA/view?usp=share_link .  Decompress the zip file and place in the .test_files directory.`)
	v33 := templates.Version33()
	r, err := os.Open("../.test_files/full-gfs-file.grb2")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()
	f := file.NewGribFile(r, v33)

	ec, err := os.Open("../.test_files/full-gfs-file.grb2.floats")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ec.Close() }()
	source := &floatReader{r: ec}

	recNum := 0
	for rec, err := range f.Records {
		recNum++
		if err != nil {
			t.Fatal(err)
		}
		def, err := rec.DataRepresentation.Definition()
		if err != nil {
			t.Fatal(err)
		}
		values, err := def.GetValues(rec)
		if err != nil {
			t.Fatal(err)
		}
		scale := def.DecimalScale()
		expectedPrecision := math.Pow(10, -float64(scale-1))
		for i, value := range values {
			expected, err := source.Next()
			if err != nil {
				t.Fatal(err)
			}
			expected64 := float64(expected)
			// we want the values to match to the desired precision. NaN values should match either NaN or eccode's placeholder value of 9999.0
			if math.Abs(value-expected64) > expectedPrecision || !(math.IsNaN(value) == math.IsNaN(expected64) || (math.IsNaN(value) && expected64 == 9999.0)) {
				t.Fatalf(`error in message %v, index %v: expected %.10f but got %.10f with expected precision %v`, recNum, i, expected, value, expectedPrecision)
			}
		}
		println(fmt.Sprintf(`finished message %v`, recNum))
	}
}
