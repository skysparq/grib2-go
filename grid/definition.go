package grid

import "github.com/skysparq/grib2-go/record"

type Definition interface {
	Parse(section record.Section3) error
}
