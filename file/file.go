package file

import (
	"errors"
	"io"

	"github.com/skysparq/grib2-go/record"
)

type GribFile interface {
	Records(yield func(record.Record, error) bool)
}

func NewGribFile(r io.Reader, template record.Templates) GribFile {
	return &gribFile{r: r, template: template}
}

type gribFile struct {
	template record.Templates
	r        io.Reader
}

func (g *gribFile) Records(yield func(record.Record, error) bool) {
	for {
		rec, err := record.ParseRecord(g.r, g.template)
		if errors.Is(err, io.EOF) {
			return
		}
		if !yield(rec, err) {
			return
		}
	}
}
