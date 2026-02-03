package utility_test

import (
	"testing"

	"github.com/skysparq/grib2-go/utility"
)

func TestStdLatLngToFloat(t *testing.T) {
	tests := []struct {
		value int
		want  float64
	}{
		{value: 8005000, want: 8.005},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := utility.StdLatLngToFloat(tt.value); got != tt.want {
				t.Errorf("StdLatLngToFloat(%d) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
