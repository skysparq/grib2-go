package projections

func normalizeLng(lng int) int {
	if lng > 180000000 {
		lng -= 360000000
	}
	return lng
}
