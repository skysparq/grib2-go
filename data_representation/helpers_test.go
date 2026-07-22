package data_representation_test

import (
	"github.com/skysparq/grib2-go/grid"
	"github.com/skysparq/grib2-go/record"
)

// fakeTemplates resolves any Section3 to a fixed GridDefinition, letting tests exercise
// GetValues without parsing a real grid definition template.
type fakeTemplates struct {
	grid record.GridDefinition
}

func (f fakeTemplates) GridDefinitionEnd(int, []byte) (int, bool) { return 0, true }

func (f fakeTemplates) ProductDefinitionEnd(int, []byte) (int, bool) { return 0, true }

func (f fakeTemplates) DataRepresentation(record.Section5) (record.DataRepresentationDefinition, error) {
	return nil, nil
}

func (f fakeTemplates) GridDefinition(record.Section3) (record.GridDefinition, error) {
	return f.grid, nil
}

func (f fakeTemplates) ProductDefinition(record.Section4) (record.ProductDefinition, error) {
	return nil, nil
}

// singleRowGridTemplates returns Templates that resolve to a 1 x totalPoints grid with the
// default (already north-up, west-first) scanning mode, so GetValues can be tested without
// needing normalization to change the result.
func singleRowGridTemplates(totalPoints int) record.Templates {
	return fakeTemplates{grid: grid.Template0{PointsAlongParallel: totalPoints, PointsAlongMeridian: 1}}
}
