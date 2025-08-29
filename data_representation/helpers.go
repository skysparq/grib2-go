package data_representation

import "math"

func getDecimalScaledRef(decimalScaleFactor int, ref float32) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * float64(ref)
}

func getScale(decimalScaleFactor int, binaryScaleFactor int) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * math.Pow(2, float64(binaryScaleFactor))
}
