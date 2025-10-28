package data_representation

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
)

type BitmapReader struct {
	bitmap     []byte
	checkIsSet func(int) bool
}

func NewBitmapReader(rec record.Record) (*BitmapReader, error) {
	r := &BitmapReader{}
	switch rec.BitMap.BitmapIndicator {
	case 0:
		r.checkIsSet = r.isSet
		r.bitmap = rec.BitMap.BitmapData
	case 255:
		r.checkIsSet = r.alwaysFalse
	default:
		return r, fmt.Errorf(`error creating BitmapReader: bitmap indicator must be 0 or 255`)
	}
	return r, nil
}

func (r *BitmapReader) IsSet(index int) bool {
	return r.checkIsSet(index)
}

func (r *BitmapReader) alwaysFalse(_ int) bool {
	return false
}

func (r *BitmapReader) isSet(index int) bool {
	return r.bitmap[index/8]&(1<<(index%8)) != 0
}
