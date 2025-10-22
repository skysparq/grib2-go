package data_representation

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
)

func checkSectionNum(section record.Section5, expectedNum int) error {
	if section.DataRepresentationTemplateNumber != expectedNum {
		return fmt.Errorf(`error parsing data representation template %d: section 5 template number is %d rather than %d`, expectedNum, section.DataRepresentationTemplateNumber, expectedNum)
	}
	return nil
}
