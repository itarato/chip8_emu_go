package chip8

func U16Mask(n uint16, m uint16) uint16 {
	for m&0xF == 0 {
		n = n >> 4
		m = m >> 4
	}
	return n & m
}