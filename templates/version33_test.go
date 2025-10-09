package templates_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/templates"
)

func TestAllFiles(t *testing.T) {
	dir, err := os.ReadDir(`../.test_files`)
	if err != nil {
		t.Fatal(err)
	}
	template := templates.Version33()
	for _, fEntry := range dir {
		if !slices.Contains([]string{`.grb2`, `.grib2`}, filepath.Ext(fEntry.Name())) {
			continue
		}
		path := filepath.Join(`../.test_files`, fEntry.Name())
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf(`error reading %v: %v`, path, err)
		}
		grib := file.NewGribFile(f, template)
		var index int
		for rec, err := range grib.Records {
			if err != nil {
				_ = f.Close()
				t.Fatalf(`error on message %v reading %v: %v`, index, path, err)
			}
			_, err = rec.Grid.Definition()
			if err != nil {
				_ = f.Close()
				t.Fatalf(`error parsing grid definition on message %v reading %v: %v`, index, path, err)
			}
			_, err = rec.Product.Definition()
			if err != nil {
				_ = f.Close()
				t.Fatalf(`error parsing product definition on message %v reading %v: %v`, index, path, err)
			}
			_, err = rec.DataRepresentation.Definition()
			if err != nil {
				_ = f.Close()
				t.Fatalf(`error parsing data representation definition on message %v reading %v: %v`, index, path, err)
			}
			index++
		}
		t.Logf(`file %v ok`, path)
		_ = f.Close()
	}
}
