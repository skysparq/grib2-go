package data_representation

func IntFromBytes(b []byte) int {
	var i int
	sign := 1
	shift := len(b) - 1
	for pos := 0; pos < len(b); pos++ {
		if pos == 0 && int(b[pos])&0x80 != 0 {
			sign = -1
			i += int(b[pos]&0x7F) << uint(shift*8)
		} else {
			i += int(b[pos]) << uint(shift*8)
		}
		shift--
	}
	return i * sign
}

func UintFromBytes(b []byte) int {
	i := 0
	shift := len(b) - 1
	for pos := 0; pos < len(b); pos++ {
		i += int(b[pos]) << uint(shift*8)
		shift--
	}
	return i
}
