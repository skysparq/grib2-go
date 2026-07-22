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
	case 2:
		// IAU 1965 oblate spheroid; treated as spherical using its major axis, consistent with the
		// other cases here since this package doesn't model ellipsoid flattening.
		return 6378160.0, nil
	case 6:
		return 6371229.0, nil
	default:
		return 0, fmt.Errorf("unsupported earth shape %d", earthShape)
	}
}

// plateCarreeSrsWkt builds the SRS WKT for an unprojected latitude/longitude grid (used by
// Template0 and Template40, which differ only in how points are spaced along a parallel). These
// grids have no standard parallel of their own - GRIB2 only gives us the corner the scan starts
// at, which is not a projection parameter.
//
// This is a plain GEOGCS with a custom PRIMEM shifted to the grid's central meridian, not a
// PROJCS/Plate_Carree wrapper: PROJCS/Plate_Carree (+proj=eqc) produces linear (meter) output by
// construction, so declaring its unit as degree is self-contradictory and confuses GDAL/PROJ's
// transform pipeline - it either errors ("inconsistent unit type between xy_in and xy_out") or
// silently returns values in meters mislabeled as degrees. A shifted PRIMEM is the WKT construct
// for "the same geographic coordinates, just measured from a different meridian": it stays
// honestly in degrees, and produces a clean monotonic longitude axis across the antimeridian
// instead of a wraparound discontinuity.
func plateCarreeSrsWkt(earthShape, firstLongitude int) (string, error) {
	radius, err := earthRadius(earthShape)
	if err != nil {
		return "", fmt.Errorf(`error generating SRS WKT: %w`, err)
	}

	centralMeridian := u.StdLatLngToFloat(u.ShiftLongitude(firstLongitude))

	return fmt.Sprintf(
		`GEOGCS["unnamed", DATUM["unknown", SPHEROID["unnamed", %.1f, 0]], PRIMEM["shifted", %v], UNIT["degree", 0.0174532925199433]]`,
		radius, centralMeridian), nil
}
