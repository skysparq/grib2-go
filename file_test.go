package grib2_go_test

import (
	grib2_go "github.com/skysparq/grib2-go"
	"github.com/skysparq/grib2-go/templates"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestLoadGribFile(t *testing.T) {
	r, err := os.Open(`./.large_test_files/full_gfs_file.grb2`)
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

	grib := grib2_go.NewGribFile(r, templates.Revision20120111())
	gridDefs := make(map[int]int)
	prodDefs := make(map[int]int)
	totalLength := 0
	for rec, recErr := range grib.Records {
		if recErr != nil {
			t.Fatal(recErr)
		}
		gridDefs[rec.GridDefinition.GridDefinitionTemplateNumber]++
		prodDefs[rec.ProductDefinition.ProductDefinitionTemplateNumber]++
		totalLength += rec.Indicator.GribLength
	}

	if expected := 538_804_943; totalLength != expected {
		t.Fatalf(`expected file size %d but got %d`, expected, totalLength)
	}
	close(done)
	t.Logf("Peak memory usage (Alloc): %d bytes (%.2f MB)", peak, float64(peak)/1024/1024)
	t.Logf(`grid definition templates: %v, product definition templates: %v`, gridDefs, prodDefs)
}
