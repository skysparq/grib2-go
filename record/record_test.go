package record_test

import (
	"fmt"
	"github.com/skysparq/grib2-go/record"
	"os"
	"testing"
)

func TestParseRecord(t *testing.T) {
	r, err := os.Open(`../test_files/single-grib2-record.grb2`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	rec, err := record.ParseRecord(r)
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
	return nil
}
