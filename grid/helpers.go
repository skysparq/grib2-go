package grid

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
	"github.com/skysparq/grib2-go/utility"
)

func checkSectionNum(section record.Section3, expectedNum int) error {
	if section.GridDefinitionTemplateNumber != expectedNum {
		return fmt.Errorf(`error parsing grid definition template %d: section 3 template number is %d rather than %d`, expectedNum, section.GridDefinitionTemplateNumber, expectedNum)
	}
	return nil
}

func mmToMeters(value int) float64 {
	return utility.ScaleInt(value, 3)
}
