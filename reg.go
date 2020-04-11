package chip8

// V[0x0F] Flag: addition, substraction, pixel draw collision.

type Reg struct {
	V  [0x10]byte
	I  uint16 // Address.
	PC uint16 // Program counter.
}

func (r *Reg) Init() {
	r.PC = 0x0200
	r.I = 0x0000
	for i, _ := range r.V {
		r.V[i] = 0x00
	}
}
