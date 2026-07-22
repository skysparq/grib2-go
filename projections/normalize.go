package projections

// NormalizeScanOrder reorders a flat Ni x Nj grid of values from GRIB2 scan order into row-major
// order with row 0 at the north edge and column 0 at the west edge, matching the orientation
// raster formats such as GeoTIFF expect. It works purely from grid dimensions and scanning mode -
// no coordinate math is involved, so the same function reorders latitudes, longitudes, and data
// values identically.
func NormalizeScanOrder(values []float64, ni, nj int, mode ScanningMode) []float64 {
	out := make([]float64, len(values))
	for row := 0; row < nj; row++ {
		j := row
		if !mode.TopToBottom {
			j = nj - 1 - row
		}
		for col := 0; col < ni; col++ {
			i := col
			if mode.RightToLeft {
				i = ni - 1 - col
			}
			var k int
			if mode.OverFirst {
				k = j*ni + i
			} else {
				k = i*nj + j
			}
			out[row*ni+col] = values[k]
		}
	}
	return out
}
