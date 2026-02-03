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
		{x: 2, n: 10, want: 1024},
		{x: 3, n: 5, want: 243},
		{x: 5, n: 0, want: 1},
		{x: 2, n: -2, want: 0.25},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := utility.PowInt(tt.x, tt.n); got != tt.want {
				t.Errorf("PowInt(%f, %d) = %v, want %v", tt.x, tt.n, got, tt.want)
			}
		})
	}
}
