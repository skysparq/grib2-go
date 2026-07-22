package grid

import (
	"fmt"

	"github.com/skysparq/grib2-go/record"
	u "github.com/skysparq/grib2-go/utility"
)

func checkSectionNum(section record.Section3, expectedNum int) error {
	if section.GridDefinitionTemplateNumber != expectedNum {
		return fmt.Errorf(`error parsing grid definition template %d: section 3 template number is %d rather than %d`, expectedNum, section.GridDefinitionTemplateNumber, expectedNum)
	}
	return nil
}

func mmToMeters(value int) float64 {
	return float64(value) * 1e-3
}

func earthRadius(earthShape int) (float64, error) {
	switch earthShape {
	case 0:
		return 6367470.0, nil
	case 6:
		return 6371229.0, nil
	default:
		return 0, fmt.Errorf("unsupported earth shape %d", earthShape)
	}
}

// plateCarreeSrsWkt builds the SRS WKT for an unprojected latitude/longitude grid (used by
// Template0 and Template40, which differ only in how points are spaced along a parallel). These
// grids have no standard parallel of their own - GRIB2 only gives us the corner the scan starts
// at, which is not a projection parameter - so the standard parallel is always 0 (true Plate
// Carree), not the grid's starting latitude.
func plateCarreeSrsWkt(earthShape, firstLongitude int) (string, error) {
	radius, err := earthRadius(earthShape)
	if err != nil {
		return "", fmt.Errorf(`error generating SRS WKT: %w`, err)
	}

	centralMeridian := u.StdLatLngToFloat(u.ShiftLongitude(firstLongitude))
	projection := fmt.Sprintf(
		`PROJECTION["Plate_Carree"], PARAMETER["False_Easting", 0], PARAMETER["False_Northing", 0], PARAMETER["Central_Meridian", %v]`,
		centralMeridian)

	return fmt.Sprintf(
		`PROJCS["unnamed", GEOGCS["unnamed", DATUM["unknown", SPHEROID["unnamed", %.1f, 0]], PRIMEM["Greenwich", 0], UNIT["degree", 0.0174532925199433]], %s, UNIT["degree", 0.0174532925199433]]`,
		radius, projection), nil
}
