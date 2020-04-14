package chip8

import (
	"image/color"
	"math/rand"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

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
	Stack    [0x10]uint16
	SP       uint8

	Win *pixelgl.Window
}

func MakeEmu(win *pixelgl.Window, ui_scale int, rom_path string) Emu {
	return Emu{
		Mem:     Mem{},
		Reg:     Reg{},
		Display: MakeDisplay(float64(ui_scale)),
		Timer:   Timer{},
		Input:   Input{},
		Sound:   Sound{Timer: Timer{}},

		Halted:   true,
		DrawFlag: false,
		RomPath:  rom_path,
		Stack:    [0x10]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		SP:       0,

		Win: win,
	}
}

func (e *Emu) Init() {
	e.Display.Init()
	e.Input.Init()
	e.Reg.Init()
	e.Mem.Init(e.RomPath)

	e.Halted = false
	e.DrawFlag = false
	e.SP = 0x00

	e.Win.Clear(colornames.Black)
}

func (e *Emu) Run() {
	e.Init()
	e.RunCycle()
}

func (e *Emu) RunCycle() {
	// f, err := os.Create("/home/itarato/Desktop/prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	for !e.Win.Closed() && !e.Halted {
		opcode, err_opcode := e.FetchOpcode()
		if err_opcode != nil {
			panic("Illegal opcode fetch")
		}

		e.ExecIntruction(opcode)

		if e.DrawFlag {
			imd := imdraw.New(nil)
			imd.Color = pixel.RGB(1, 1, 1)

			e.Win.Clear(color.Black)
			e.Display.Draw(imd)
			imd.Draw(e.Win)
			e.Win.Update()
			e.DrawFlag = false
		}

		e.Input.UpdateState(e.Win)
		e.Sound.Update()
		e.Timer.Dec() // @TODO - Probably needs ignoring opcode while in-delay

		// @TODO Throttle 1 cycle to 16ms (60 instruction per second).
	}
}

func (e *Emu) ExecIntruction(opcode uint16) {
	reg_num_x := U16Mask(opcode, 0x0F00)
	reg_num_y := U16Mask(opcode, 0x00F0)

	switch U16Mask(opcode, 0xF000) {
	case 0x0:
		switch U16Mask(opcode, 0x00FF) {
		case 0xE0:
			// 00E0 	Display 	disp_clear() 	Clears the screen.
			e.Display.Clear()
			e.DrawFlag = true
		case 0xEE:
			// 00EE 	Flow 	return; 	Returns from a subroutine.
			e.Reg.PC = e.StackPop()
		default:
			// 0NNN 	Call 		Calls RCA 1802 program at address NNN. Not necessary for most ROMs.
			panic("Unimplemented opcode 0NNN")
		}
	case 0x1:
		// 1NNN 	Flow 	goto NNN; 	Jumps to address NNN.
		e.Reg.PC = U16Mask(opcode, 0x0FFF)
	case 0x2:
		// 2NNN 	Flow 	*(0xNNN)() 	Calls subroutine at NNN.
		e.StackPush(e.Reg.PC)
		e.Reg.PC = U16Mask(opcode, 0x0FFF)
	case 0x3:
		// 3XNN 	Cond 	if(Vx==NN) 	Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)
		nn := U16Mask(opcode, 0x00FF)
		if e.Reg.V[reg_num_x] == uint8(nn) {
			e.Reg.SkipOpcode()
		}
	case 0x4:
		// 4XNN 	Cond 	if(Vx!=NN) 	Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
		nn := U16Mask(opcode, 0x00FF)
		if e.Reg.V[reg_num_x] != uint8(nn) {
			e.Reg.SkipOpcode()
		}
	case 0x5:
		// 5XY0 	Cond 	if(Vx==Vy) 	Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)
		if e.Reg.V[reg_num_x] == e.Reg.V[reg_num_y] {
			e.Reg.SkipOpcode()
		}
	case 0x6:
		// 6XNN 	Const 	Vx = NN 	Sets VX to NN.
		nn := U16Mask(opcode, 0x00FF)
		e.Reg.V[reg_num_x] = uint8(nn)
	case 0x7:
		// 7XNN 	Const 	Vx += NN 	Adds NN to VX. (Carry flag is not changed)
		nn := U16Mask(opcode, 0x00FF)
		e.Reg.V[reg_num_x] += uint8(nn)
	case 0x8:
		switch U16Mask(opcode, 0x000F) {
		case 0x0:
			// 8XY0 	Assign 	Vx=Vy 	Sets VX to the value of VY.
			e.Reg.V[reg_num_x] = e.Reg.V[reg_num_y]
		case 0x1:
			// 8XY1 	BitOp 	Vx=Vx|Vy 	Sets VX to VX or VY. (Bitwise OR operation)
			e.Reg.V[reg_num_x] |= e.Reg.V[reg_num_y]
		case 0x2:
			// 8XY2 	BitOp 	Vx=Vx&Vy 	Sets VX to VX and VY. (Bitwise AND operation)
			e.Reg.V[reg_num_x] &= e.Reg.V[reg_num_y]
		case 0x3:
			// 8XY3[a] 	BitOp 	Vx=Vx^Vy 	Sets VX to VX xor VY.
			e.Reg.V[reg_num_x] ^= e.Reg.V[reg_num_y]
		case 0x4:
			// 8XY4 	Math 	Vx += Vy 	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
			if 0xFF-e.Reg.V[reg_num_x] < e.Reg.V[reg_num_y] {
				e.Reg.SetRegVF()
			} else {
				e.Reg.UnsetRegVF()
			}
			e.Reg.V[reg_num_x] += e.Reg.V[reg_num_y]
		case 0x5:
			// 8XY5 	Math 	Vx -= Vy 	VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
			if e.Reg.V[reg_num_x] < e.Reg.V[reg_num_y] {
				e.Reg.UnsetRegVF()
			} else {
				e.Reg.SetRegVF()
			}
			e.Reg.V[reg_num_x] -= e.Reg.V[reg_num_y]
		case 0x6:
			// 8XY6[a] 	BitOp 	Vx>>=1 	Stores the least significant bit of VX in VF and then shifts VX to the right by 1.[b]
			e.Reg.V[0xF] = e.Reg.V[reg_num_x] & 1
			e.Reg.V[reg_num_x] >>= 1
		case 0x7:
			// 8XY7[a] 	Math 	Vx=Vy-Vx 	Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
			if e.Reg.V[reg_num_x] > e.Reg.V[reg_num_y] {
				e.Reg.UnsetRegVF()
			} else {
				e.Reg.SetRegVF()
			}
			e.Reg.V[reg_num_x] = e.Reg.V[reg_num_y] - e.Reg.V[reg_num_x]
		case 0xE:
			// 8XYE[a] 	BitOp 	Vx<<=1 	Stores the most significant bit of VX in VF and then shifts VX to the left by 1.[b]
			e.Reg.V[0xF] = (e.Reg.V[reg_num_x] >> 7) & 1
			e.Reg.V[reg_num_x] <<= 1
		default:
			panic("Illegal 8xxx opcode")
		}
	case 0x9:
		// 9XY0 	Cond 	if(Vx!=Vy) 	Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
		if e.Reg.V[reg_num_x] != e.Reg.V[reg_num_y] {
			e.Reg.SkipOpcode()
		}
	case 0xA:
		// ANNN 	MEM 	I = NNN 	Sets I to the address NNN.
		e.Reg.I = U16Mask(opcode, 0x0FFF)
	case 0xB:
		// BNNN 	Flow 	PC=V0+NNN 	Jumps to the address NNN plus V0.
		e.Reg.PC = U16Mask(opcode, 0x0FFF) + uint16(e.Reg.V[0])
	case 0xC:
		// CXNN 	Rand 	Vx=rand()&NN 	Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
		rand := uint8(rand.Uint32() & 0xFF)
		mask := uint8(U16Mask(opcode, 0x00FF))
		e.Reg.V[reg_num_x] = rand & mask
	case 0xD:
		// DXYN 	Disp 	draw(Vx,Vy,N) 	Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.
		// 				Each row of 8 pixels is read as bit-coded starting from memory location I; I value doesn’t change after the execution of this instruction.
		// 				As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that doesn’t happen
		n := U16Mask(opcode, 0x000F)
		byte_arr, err_byte_arr := e.Mem.ReadRange(e.Reg.I, n)
		if err_byte_arr != nil {
			panic("Cannot read sprite data")
		}

		has_flip := e.Display.DrawSprite(uint8(e.Reg.V[reg_num_x]), uint8(e.Reg.V[reg_num_y]), uint8(n), byte_arr)
		if has_flip {
			e.Reg.SetRegVF()
		} else {
			e.Reg.UnsetRegVF()
		}

		e.DrawFlag = true
	case 0xE:
		switch U16Mask(opcode, 0x00FF) {
		case 0x9E:
			// EX9E 	KeyOp 	if(key()==Vx) 	Skips the next instruction if the key stored in VX is pressed.
			// (Usually the next instruction is a jump to skip a code block)
			if e.Input.IsPressed(e.Reg.V[reg_num_x]) {
				e.Reg.SkipOpcode()
			}
		case 0xA1:
			// EXA1 	KeyOp 	if(key()!=Vx) 	Skips the next instruction if the key stored in VX isn't pressed.
			// (Usually the next instruction is a jump to skip a code block)
			if !e.Input.IsPressed(e.Reg.V[reg_num_x]) {
				e.Reg.SkipOpcode()
			}
		default:
			panic("Illegal Exxx opcode")
		}
	case 0xF:
		switch U16Mask(opcode, 0x00FF) {
		case 0x07:
			// FX07 	Timer 	Vx = get_delay() 	Sets VX to the value of the delay timer.
			e.Reg.V[reg_num_x] = uint8(e.Timer.Counter)
		case 0x0A:
			// FX0A 	KeyOp 	Vx = get_key() 	A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
			panic("Unimplemented FX0A")
		case 0x15:
			// FX15 	Timer 	delay_timer(Vx) 	Sets the delay timer to VX.
			e.Timer.Set(e.Reg.V[reg_num_x])
		case 0x18:
			// FX18 	Sound 	sound_timer(Vx) 	Sets the sound timer to VX.
			e.Sound.Timer.Set(e.Reg.V[reg_num_x])
		case 0x1E:
			// FX1E 	MEM 	I +=Vx 	Adds VX to I. VF is set to 1 when there is a range overflow (I+VX>0xFFF), and to 0 when there isn't.[c]
			if e.Reg.I > 0x0FFF-uint16(e.Reg.V[reg_num_x]) {
				e.Reg.SetRegVF()
			} else {
				e.Reg.UnsetRegVF()
			}
			e.Reg.I += uint16(e.Reg.V[reg_num_x])
		case 0x29:
			// FX29 	MEM 	I=sprite_addr[Vx] 	Sets I to the location of the sprite for the character in VX.
			// 	Characters 0-F (in hexadecimal) are represented by a 4x5 font.
			font_idx := e.Reg.V[reg_num_x]
			if font_idx >= 0xF {
				panic("Unknown font index")
			}
			e.Reg.I = uint16(font_idx) * 5
		case 0x33:
			// FX33 	BCD 	set_BCD(Vx);
			// 	*(I+0)=BCD(3);
			// 	*(I+1)=BCD(2);
			// 	*(I+2)=BCD(1);
			// 	Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I,
			// 	the middle digit at I plus 1, and the least significant digit at I plus 2. (In other words, take the decimal representation of VX,
			// 	place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)
			num := e.Reg.V[reg_num_x]
			for i := 0; i < 3; i++ {
				err_write := e.Mem.Write(e.Reg.I+2-uint16(i), num%10)
				num /= 10
				if err_write != nil {
					panic("Cannot write BCD")
				}
			}
		case 0x55:
			// FX55 	MEM 	reg_dump(Vx,&I) 	Stores V0 to VX (including VX) in memory starting at address I.
			// 	The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
			reg_lim := reg_num_x
			for i := 0; i <= int(reg_lim); i++ {
				e.Mem.Write(e.Reg.I+uint16(i), e.Reg.V[i])
			}
		case 0x65:
			// FX65 	MEM 	reg_load(Vx,&I) 	Fills V0 to VX (including VX) with values from memory starting at address I.
			// 	The offset from I is increased by 1 for each value written, but I itself is left unmodified.[d]
			reg_lim := reg_num_x
			for i := 0; i <= int(reg_lim); i++ {
				b, err_b := e.Mem.Read(e.Reg.I + uint16(i))
				if err_b != nil {
					panic("Illegal read for FX65")
				}
				e.Reg.V[i] = b
			}
		default:
			panic("Illegal Fxxx opcode.")
		}
	default:
		panic("Illegal highes byte of opcode.")
	}
}

func (e *Emu) StackPush(v uint16) {
	if e.SP >= 0x10 {
		panic("Stack is full")
	}

	e.Stack[e.SP] = v
	e.SP += 1
}

func (e *Emu) StackPop() uint16 {
	if e.SP == 0 {
		panic("Stack is empty")
	}

	e.SP -= 1
	return e.Stack[e.SP]
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
