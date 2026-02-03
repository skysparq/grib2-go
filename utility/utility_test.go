package utility_test

import (
	"testing"

	"github.com/skysparq/grib2-go/utility"
)

func TestPowInt(t *testing.T) {
	tests := []struct {
		x    float64
		n    int
		want float64
	}{
		{2, 10, 1024},
		{3, 5, 243},
		{5, 0, 1},
		{2, -2, 0.25},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := utility.PowInt(tt.x, tt.n); got != tt.want {
				t.Errorf("PowInt(%f, %d) = %v, want %v", tt.x, tt.n, got, tt.want)
			}
		})
	}
}
