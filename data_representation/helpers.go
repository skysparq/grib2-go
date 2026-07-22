package data_representation

import (
	"fmt"

	"github.com/skysparq/grib2-go/projections"
	"github.com/skysparq/grib2-go/record"
)

func checkSectionNum(section record.Section5, expectedNum int) error {
	if section.DataRepresentationTemplateNumber != expectedNum {
		return fmt.Errorf(`error parsing data representation template %d: section 5 template number is %d rather than %d`, expectedNum, section.DataRepresentationTemplateNumber, expectedNum)
	}
	return nil
}

// normalizeScanOrder reorders decoded values from GRIB2 scan order into row-major order with
// north at row 0 and west at column 0, matching the order grid.GridDefinition.Points() uses.
func normalizeScanOrder(rec record.Record, values []float64) ([]float64, error) {
	grid, err := rec.Grid.Definition()
	if err != nil {
		return nil, fmt.Errorf("error normalizing scan order: %w", err)
	}
	return projections.NormalizeScanOrder(values, grid.XVals(), grid.YVals(), grid.ScanMode()), nil
}
