package record_test

import (
	"fmt"
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

func TestSamplePointsFromFullGrib(t *testing.T) {
	v33 := templates.Version33()
	_, r, err := test_files.Load(test_files.FullGfsFile)
	if err != nil {
		t.Fatal(err)
	}
	w, err := os.Create(`../.test_files/full_gfs_sample.txt`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = w.Close() }()

	f := file.NewGribFile(r, v33)
	recNum := 0
	for rec, err := range f.Records {
		recNum++
		if err != nil {
			t.Fatal(err)
		}
		vals, err := rec.GetGriddedValues()
		if err != nil {
			t.Fatalf("error on rec %v: %v", recNum, err.Error())
		}
		for y := 0; y < vals.YVals; y += 200 {
			for x := 0; x < vals.XVals; x += 200 {
				index := y*vals.XVals + x
				lng := vals.Lngs[index]
				lat := vals.Lats[index]
				val := vals.Values[index]
				_, err = w.WriteString(fmt.Sprintf("%v,%f,%f,%f\n", recNum, lng, lat, val))
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}

}
