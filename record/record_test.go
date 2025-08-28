package record_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/skysparq/grib2/record"
	"github.com/skysparq/grib2/templates"
)

func TestParseGfsRecordPointInTime(t *testing.T) {
	template := templates.Revision20120111()
	r, err := os.Open(`../test_files/single-grib2-record.grb2`)
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
	err = checkSection3(rec.GridDefinition)
	if err != nil {
		t.Fatal(err)
	}
	err = checkSection4(rec.ProductDefinition)
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
	if expected := 21; section.Length != expected {
		return fmt.Errorf(`expected length %d, got %d`, expected, section.Length)
	}
	if expected := 7; section.OriginatingCenter != expected {
		return fmt.Errorf(`expected originating center %d, got %d`, expected, section.OriginatingCenter)
	}
	if expected := 0; section.OriginatingSubCenter != expected {
		return fmt.Errorf(`expected originating subcenter %d, got %d`, expected, section.OriginatingSubCenter)
	}
	if expected := 2; section.MasterTableVersion != expected {
		return fmt.Errorf(`expected master table version %d, got %d`, expected, section.MasterTableVersion)
	}
	if expected := 1; section.LocalTableVersion != expected {
		return fmt.Errorf(`expected local table version %d, got %d`, expected, section.LocalTableVersion)
	}
	if expected := 1; section.ReferenceTimeSignificance != expected {
		return fmt.Errorf(`expected reference time significance %d, got %d`, expected, section.ReferenceTimeSignificance)
	}
	if expected := 2025; section.Year != expected {
		return fmt.Errorf(`expected year %d, got %d`, expected, section.Year)
	}
	if expected := 3; section.Month != expected {
		return fmt.Errorf(`expected month %d, got %d`, expected, section.Month)
	}
	if expected := 5; section.Day != expected {
		return fmt.Errorf(`expected day %d, got %d`, expected, section.Day)
	}
	if expected := 6; section.Hour != expected {
		return fmt.Errorf(`expected hour %d, got %d`, expected, section.Hour)
	}
	if expected := 0; section.Minute != expected {
		return fmt.Errorf(`expected minute %d, got %d`, expected, section.Minute)
	}
	if expected := 0; section.Second != expected {
		return fmt.Errorf(`expected second %d, got %d`, expected, section.Second)
	}
	if expected := 0; section.ProductionStatus != expected {
		return fmt.Errorf(`expected production status %d, got %d`, expected, section.ProductionStatus)
	}
	if expected := 1; section.DataType != expected {
		return fmt.Errorf(`expected data type %d, got %d`, expected, section.DataType)
	}
	if expected := 0; len(section.Reserved) != expected {
		return fmt.Errorf(`expected reserved length %d, got %d`, expected, len(section.Reserved))
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
	if expected := 72; section.Length != expected {
		return fmt.Errorf(`expected length %d, got %d`, expected, section.Length)
	}
	if expected := 0; section.GridSourceDefinition != expected {
		return fmt.Errorf(`expected grid source definition %d, got %d`, expected, section.GridSourceDefinition)
	}
	if expected := 1_038_240; section.TotalPoints != expected {
		return fmt.Errorf(`expected total points %d, got %d`, expected, section.TotalPoints)
	}
	if expected := 0; section.OctetsForOptionalPointList != expected {
		return fmt.Errorf(`expected octets for optional point list %d, got %d`, expected, section.OctetsForOptionalPointList)
	}
	if expected := 0; section.InterpretationOfPointList != expected {
		return fmt.Errorf(`expected interpretation of point list %d, got %d`, expected, section.InterpretationOfPointList)
	}
	if expected := 72 - 15 + 1; len(section.GridDefinitionTemplateData) != expected {
		return fmt.Errorf(`expected grid definition template length %d, got %d`, expected, len(section.GridDefinitionTemplateData))
	}
	if expected := 0; len(section.OptionalPointListData) != expected {
		return fmt.Errorf(`expectedg optional point list length %d, got %d`, expected, len(section.OptionalPointListData))
	}
	return nil
}

func checkSection4(section record.Section4) error {
	if expected := 34; section.Length != expected {
		return fmt.Errorf(`expected length %d, got %d`, expected, section.Length)
	}
	if expected := 0; section.CoordinateValuesAfterTemplate != expected {
		return fmt.Errorf(`expected coordinate values after template %d, got %d`, expected, section.CoordinateValuesAfterTemplate)
	}
	if expected := 0; section.ProductDefinitionTemplateNumber != expected {
		return fmt.Errorf(`expected product template number %d, got %d`, expected, section.ProductDefinitionTemplateNumber)
	}
	if expected := 25; len(section.ProductDefinitionTemplateData) != expected {
		return fmt.Errorf(`expected product template length %d, got %d`, expected, len(section.ProductDefinitionTemplateData))
	}
	if expected := 0; len(section.CoordinateValuesData) != expected {
		return fmt.Errorf(`expected coordinate values data length %d, got %d`, expected, len(section.CoordinateValuesData))
	}
	return nil
}

func checkSection5(section record.Section5) error {
	if expected := 49; section.Length != expected {
		return fmt.Errorf(`expected length %d, got %d`, expected, section.Length)
	}
	if expected := 1_038_240; section.TotalDataPoints != expected {
		return fmt.Errorf(`expected coordinate values after template %d, got %d`, expected, section.TotalDataPoints)
	}
	if expected := 3; section.DataRepresentationTemplateNumber != expected {
		return fmt.Errorf(`expected data representation template %d, got %d`, expected, section.DataRepresentationTemplateNumber)
	}
	if expected := 38; len(section.DataRepresentationTemplateData) != expected {
		return fmt.Errorf(`expected data representation length %d, got %d`, expected, len(section.DataRepresentationTemplateData))
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
	template := templates.Revision20120111()
	r, err := os.Open(`../test_files/single-grib2-record-prod-def-8.grb2`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	_, err = record.ParseRecord(r, template)
	if err != nil {
		t.Fatal(err)
	}
}
