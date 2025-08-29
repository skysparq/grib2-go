package data_representation

import (
	"fmt"
	"math"

	"github.com/skysparq/grib2-go/record"
)

func checkSectionNum(section record.Section5, expectedNum int) error {
	if section.DataRepresentationTemplateNumber != expectedNum {
		return fmt.Errorf(`error parsing data representation template %d: section 3 template number is %d rather than %d`, expectedNum, section.DataRepresentationTemplateNumber, expectedNum)
	}
	return nil
}

func getDecimalScaledRef(decimalScaleFactor int, ref float32) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * float64(ref)
}

func getScale(decimalScaleFactor int, binaryScaleFactor int) float64 {
	return math.Pow(10, -float64(decimalScaleFactor)) * math.Pow(2, float64(binaryScaleFactor))
}
