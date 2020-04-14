package chip8

type Reg struct {
	V  [0x10]uint8
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

func (r *Reg) SkipOpcode() {
	r.PC += 2
}

func (r *Reg) SetRegVF() {
	r.V[0xF] = 1
}

func (r *Reg) UnsetRegVF() {
	r.V[0xF] = 0
}
