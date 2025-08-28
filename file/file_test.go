package file_test

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/templates"
)

func TestLoadGribFile(t *testing.T) {
	path := `../.large_test_files/full_gfs_file.grb2`
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	r, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = r.Close() }()

	var peak uint64
	done := make(chan struct{})
	go func() {
		var m runtime.MemStats
		for {
			select {
			case <-done:
				return
			default:
				runtime.ReadMemStats(&m)
				if m.Alloc > peak {
					peak = m.Alloc
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	grib := file.NewGribFile(r, templates.Revision20120111())
	gridDefs := make(map[int]int)
	prodDefs := make(map[int]int)
	dataRepDefs := make(map[int]int)
	totalLength := 0
	for rec, recErr := range grib.Records {
		if recErr != nil {
			t.Fatal(recErr)
		}
		gridDefs[rec.GridDefinition.GridDefinitionTemplateNumber]++
		prodDefs[rec.ProductDefinition.ProductDefinitionTemplateNumber]++
		dataRepDefs[rec.DataRepresentation.DataRepresentationTemplateNumber]++
		totalLength += rec.Indicator.GribLength
	}

	if expected := int(stat.Size()); totalLength != expected {
		t.Fatalf(`expected file size %d but got %d`, expected, totalLength)
	}
	close(done)
	t.Logf("Peak memory usage (Alloc): %d bytes (%.2f MB)", peak, float64(peak)/1024/1024)
	t.Logf(`grid definition templates: %v, product definition templates: %v, data representation templates: %v`, gridDefs, prodDefs, dataRepDefs)
}
