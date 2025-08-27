package grib2_go

import (
	"errors"
	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/templates"
	"io"
)

type GribFile interface {
	Records(yield func(record.Record, error) bool)
}

func NewGribFile(r io.Reader, template templates.Template) GribFile {
	return &gribFile{r: r, template: template}
}

type gribFile struct {
	template templates.Template
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
