package chip8

const (
	ModeChip8 = iota
)

type Emu struct {
	Mem
	Reg
	Display
	Sound
	Timer
	Input
	Halted   bool
	DrawFlag bool
	RomPath  string
}

func MakeEmu(rom_path string) Emu {
	return Emu{
		Mem:     Mem{},
		Reg:     Reg{},
		Display: Display{},
		Sound: Sound{
			Timer: Timer{},
		},
		Timer:    Timer{},
		Input:    Input{},
		Halted:   true,
		DrawFlag: false,
		RomPath:  rom_path,
	}
}

func (e *Emu) Init() {
	e.Display.Init()
	e.Input.Init()
	e.Reg.Init()
	e.Mem.Init(e.RomPath)

	e.Halted = false
	e.DrawFlag = false
}

func (e *Emu) Run() {
	e.Init()
	e.RunCycle()
}

func (e *Emu) RunCycle() {
	for !e.Halted {
		opcode, err_opcode := e.FetchOpcode()
		if err_opcode != nil {
			panic("Illegal opcode fetch")
		}

		e.ExecIntruction(opcode)

		if e.DrawFlag {
			e.Display.Draw()
		}

		e.Input.UpdateState()
		e.Sound.Update()
		e.Timer.Dec()

		// @TODO Throttle 1 cycle to 16ms (60 instruction per second).
	}
}

func (e *Emu) ExecIntruction(opcode uint16) {
	switch U16Mask(opcode, 0xF000) {
	case 0x0:
	// 0NNN 	Call 		Calls RCA 1802 program at address NNN. Not necessary for most ROMs.
	// 00E0 	Display 	disp_clear() 	Clears the screen.
	// 00EE 	Flow 	return; 	Returns from a subroutine.
	case 0x1:
	// 1NNN 	Flow 	goto NNN; 	Jumps to address NNN.
	case 0x2:
	// 2NNN 	Flow 	*(0xNNN)() 	Calls subroutine at NNN.
	case 0x3:
	// 3XNN 	Cond 	if(Vx==NN) 	Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)
	case 0x4:
	// 4XNN 	Cond 	if(Vx!=NN) 	Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
	case 0x5:
	// 5XY0 	Cond 	if(Vx==Vy) 	Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)
	case 0x6:
	// 6XNN 	Const 	Vx = NN 	Sets VX to NN.
	case 0x7:
	// 7XNN 	Const 	Vx += NN 	Adds NN to VX. (Carry flag is not changed)
	case 0x8:
	// 8XY0 	Assign 	Vx=Vy 	Sets VX to the value of VY.
	// 8XY1 	BitOp 	Vx=Vx|Vy 	Sets VX to VX or VY. (Bitwise OR operation)
	// 8XY2 	BitOp 	Vx=Vx&Vy 	Sets VX to VX and VY. (Bitwise AND operation)
	// 8XY3[a] 	BitOp 	Vx=Vx^Vy 	Sets VX to VX xor VY.
	// 8XY4 	Math 	Vx += Vy 	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
	// 8XY5 	Math 	Vx -= Vy 	VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
	// 8XY6[a] 	BitOp 	Vx>>=1 	Stores the least significant bit of VX in VF and then shifts VX to the right by 1.[b]
	// 8XY7[a] 	Math 	Vx=Vy-Vx 	Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
	// 8XYE[a] 	BitOp 	Vx<<=1 	Stores the most significant bit of VX in VF and then shifts VX to the left by 1.[b]
	case 0x9:
	// 9XY0 	Cond 	if(Vx!=Vy) 	Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
	case 0xA:
	// ANNN 	MEM 	I = NNN 	Sets I to the address NNN.
	case 0xB:
	// BNNN 	Flow 	PC=V0+NNN 	Jumps to the address NNN plus V0.
	case 0xC:
	// CXNN 	Rand 	Vx=rand()&NN 	Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
	case 0xD:
	// DXYN 	Disp 	draw(Vx,Vy,N) 	Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels. Each row of 8 pixels is read as bit-coded starting from memory location I; I value doesn’t change after the execution of this instruction. As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that doesn’t happen
	case 0xE:
	// EX9E 	KeyOp 	if(key()==Vx) 	Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)
	// EXA1 	KeyOp 	if(key()!=Vx) 	Skips the next instruction if the key stored in VX isn't pressed. (Usually the next instruction is a jump to skip a code block)
	case 0xF:
	// FX07 	Timer 	Vx = get_delay() 	Sets VX to the value of the delay timer.
	// FX0A 	KeyOp 	Vx = get_key() 	A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
	// FX15 	Timer 	delay_timer(Vx) 	Sets the delay timer to VX.
	// FX18 	Sound 	sound_timer(Vx) 	Sets the sound timer to VX.
	// FX1E 	MEM 	I +=Vx 	Adds VX to I. VF is set to 1 when there is a range overflow (I+VX>0xFFF), and to 0 when there isn't.[c]
	// FX29 	MEM 	I=sprite_addr[Vx] 	Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
	// FX33 	BCD 	set_BCD(Vx);
	// *(I+0)=BCD(3);
	// *(I+1)=BCD(2);
	// *(I+2)=BCD(1);
	// 	Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)
	// FX55 	MEM 	reg_dump(Vx,&I) 	Stores V0 to VX (including VX) in memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
	// FX65 	MEM 	reg_load(Vx,&I) 	Fills V0 to VX (including VX) with values from memory starting at address I. The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
	default:
		panic("Illegal highes byte of opcode.")
	}
}

func (e *Emu) FetchOpcode() (uint16, error) {
	hi, err_hi := e.Mem.Read(e.Reg.PC)
	if err_hi != nil {
		return 0, err_hi
	}

	lo, err_lo := e.Mem.Read(e.Reg.PC + 1)
	if err_lo != nil {
		return 0, err_lo
	}

	e.Reg.PC += 2

	opcode := uint16(hi)<<8 | uint16(lo)
	return opcode, nil
}
