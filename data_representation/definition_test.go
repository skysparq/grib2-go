package data_representation_test

import (
	"os"
	"testing"

	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
)

func TestParseDefinitionOnGfsGrib(t *testing.T) {
	r, err := os.Open(`../.large_test_files/full_gfs_file.grb2`)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	parser := &grid.Parser{}
	grib := file.NewGribFile(r, templates.Revision20120111())
	var rec record.Record
	for rec, err = range grib.Records {
		if err != nil {
			t.Fatal(err)
		}
		_, err = parser.ParseDefinition(rec.GridDefinition)
		if err != nil {
			t.Fatal(err)
		}
	}
}
