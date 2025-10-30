package data_representation

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
)

// BitmapReader is used to read the bitmap section of a GRIB record. The bitmap section is used to determine which
// values are missing.
type BitmapReader struct {
	bitmap     []byte
	checkIsSet func(int) bool
}

// NewBitmapReader instantiates a new BitmapReader from the given record.
func NewBitmapReader(rec record.Record) (*BitmapReader, error) {
	r := &BitmapReader{}
	switch rec.BitMap.BitmapIndicator {
	case 0:
		r.checkIsSet = r.isMissing
		r.bitmap = rec.BitMap.BitmapData
	case 255:
		r.checkIsSet = r.alwaysFalse
	default:
		return r, fmt.Errorf(`error creating BitmapReader: bitmap indicator must be 0 or 255`)
	}
	return r, nil
}

// IsMissing returns true if the value at the given index is missing.
func (r *BitmapReader) IsMissing(index int) bool {
	return r.checkIsSet(index)
}

func (r *BitmapReader) alwaysFalse(_ int) bool {
	return false
}

func (r *BitmapReader) isMissing(index int) bool {
	return (r.bitmap[index/8]>>(7-(index%8)))&1 == 0
}
