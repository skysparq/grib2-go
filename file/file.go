package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/skysparq/grib2-go/record"
)

// GribFile provides an iterator over the records in a GRIB file.
type GribFile interface {
	Records(yield func(record.Record, error) bool)
}

// NewGribFile instantiates a GribFile from an io.Reader and a record.Templates.
// Standard templates can be accessed via the templates.Version33 function.
func NewGribFile(r io.Reader, template record.Templates) GribFile {
	return &gribFile{r: r, template: template}
}

type gribFile struct {
	template record.Templates
	r        io.Reader
}

// Records iterates over the records in a GRIB file.
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

// ExtractRecordBytes extracts the full blob of a GRIB record from a GRIB file.
// This is useful in isolating a single record from a GRIB file for additional testing and validation.
func ExtractRecordBytes(r io.ReadSeeker, messageNumber int) ([]byte, error) {
	for i := 0; i < messageNumber-1; i++ {
		sec0, err := record.ParseSection0(r)
		if err != nil {
			return nil, fmt.Errorf(`error extracting record %d: %w`, messageNumber, err)
		}
		_, err = r.Seek(int64(sec0.GribLength-16), io.SeekCurrent)
		if err != nil {
			return nil, fmt.Errorf(`error extracting record %d: %w`, messageNumber, err)
		}
	}
	sec0Bytes := make([]byte, 16)
	_, err := io.ReadFull(r, sec0Bytes)
	if err != nil {
		return nil, fmt.Errorf(`error extracting record %d: %w`, messageNumber, err)
	}
	sec0, err := record.ParseSection0(bytes.NewReader(sec0Bytes))
	if err != nil {
		return nil, fmt.Errorf(`error extracting record %d: %w`, messageNumber, err)
	}
	remainingLen := sec0.GribLength - 16
	remainingBytes := make([]byte, remainingLen)
	_, err = io.ReadFull(r, remainingBytes)
	if err != nil {
		return nil, fmt.Errorf(`error extracting record %d: %w`, messageNumber, err)
	}
	return append(sec0Bytes, remainingBytes...), nil
}
