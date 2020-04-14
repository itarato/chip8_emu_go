package chip8

func U16Mask(n uint16, m uint16) uint16 {
	for m&0xF == 0 {
		n >>= 4
		m >>= 4
	}
	return n & m
}

func BitN(b uint8, pos uint8) uint8 {
	return (b >> (7 - pos)) & 1
}
