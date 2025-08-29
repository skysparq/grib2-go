package grid_test

import (
	"testing"

	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestParseDefinitionOnGfsGrib(t *testing.T) {
	_, r, err := test_files.Load(test_files.FullGfsFile)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	parser := &grid.Parser{}
	grib := file.NewGribFile(r, templates.Version33())
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
