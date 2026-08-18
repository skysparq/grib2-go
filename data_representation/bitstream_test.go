package data_representation_test

import (
	"testing"

	"github.com/skysparq/grib2-go/data_representation"
)

func TestBitStream_ReadBits(t *testing.T) {
	data := []byte{0x12, 0x34, 0x56, 0x78}

	tests := []struct {
		pos  int
		n    int
		want uint64
	}{
		{0, 16, 0x1234},
		{8, 16, 0x3456},
		{16, 16, 0x5678},
		{4, 16, 0x2345},
		{0, 4, 0x1},
		{4, 4, 0x2},
		{12, 4, 0x4},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			b := data_representation.NewBitStream(data)
			b.SetPos(tt.pos)
			if got := b.ReadBits(tt.n); got != tt.want {
				t.Errorf("ReadBits(%d) from pos %d = %v, want %v", tt.n, tt.pos, got, tt.want)
			}
		})
	}
}
