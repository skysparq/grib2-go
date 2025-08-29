package product

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
)

func checkSectionNum(section record.Section4, expectedNum int) error {
	if section.ProductDefinitionTemplateNumber != expectedNum {
		return fmt.Errorf(`error parsing product definition template %d: section 4 template number is %d rather than %d`, expectedNum, section.ProductDefinitionTemplateNumber, expectedNum)
	}
	return nil
}
