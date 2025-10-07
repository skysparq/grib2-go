package file_test

import (
	"runtime"
	"testing"
	"time"

	"github.com/skysparq/grib2-go/file"
	"github.com/skysparq/grib2-go/templates"
	"github.com/skysparq/grib2-go/test_files"
)

func TestLoadGribFile(t *testing.T) {
	size, r, err := test_files.Load(test_files.FullGfsFile)
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

	grib := file.NewGribFile(r, templates.Version33())
	gridDefs := make(map[int]int)
	prodDefs := make(map[int]int)
	dataRepDefs := make(map[int]int)
	totalLength := 0
	for rec, recErr := range grib.Records {
		if recErr != nil {
			t.Fatal(recErr)
		}
		gridDefs[rec.Grid.GridDefinitionTemplateNumber]++
		prodDefs[rec.Product.ProductDefinitionTemplateNumber]++
		dataRepDefs[rec.DataRepresentation.DataRepresentationTemplateNumber]++
		totalLength += rec.Indicator.GribLength
	}

	if totalLength != size {
		t.Fatalf(`expected file size %d but got %d`, size, totalLength)
	}
	close(done)
	t.Logf("Peak memory usage (Alloc): %d bytes (%.2f MB)", peak, float64(peak)/1024/1024)
	t.Logf(`grid definition templates: %v, product definition templates: %v, data representation templates: %v`, gridDefs, prodDefs, dataRepDefs)
}
