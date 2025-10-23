package projections

import "math"

func normalizeLng(lng int) int {
	if lng > 180000000 {
		lng -= 360000000
	}
	return lng
}

// Helper functions for degree/radian conversion
func degToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func radToDeg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}
