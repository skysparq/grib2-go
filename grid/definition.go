package grid

import (
	"fmt"

	"github.com/skysparq/grib2/record"
)

type Definition interface {
	Parse(section record.Section3) (Definition, error)
}

func ParseDefinition(section record.Section3) (Definition, error) {
	switch section.GridDefinitionTemplateNumber {
	case 0:
		return (&Template0{}).Parse(section)
	default:
		return nil, fmt.Errorf(`error parsing grid definition: unsupported template number %d`, section.GridDefinitionTemplateNumber)
	}
}
