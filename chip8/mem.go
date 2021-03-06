package chip8

import (
	"fmt"
	"io/ioutil"
)

// 0x0000 - 0x01FF: ROM (mostly original version, new versions has ROM outside of 4K bank)
// +-> 0x0050 - 0x00A0: 4x5 pixel fonts.
// 0x0200 -
// 0x0EA0 - 0x0EFF: Call stack, internal use, misc vars.
// 0x0F00 - 0x0FFF: Display refresh

type Mem struct {
	Data [0x1000]uint8
}

func (m *Mem) Init(rom_path string) {
	for i, _ := range m.Data {
		m.Data[i] = 0x00
	}

	fontset := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	for i, f := range fontset {
		m.Data[i] = f
	}

	// Load ROM.
	rom_content, err_rom_content := ioutil.ReadFile(rom_path)
	if err_rom_content != nil {
		panic("Cannot read rom")
	}
	for i, b := range rom_content {
		m.Data[0x0200+i] = b
	}
}

func (m *Mem) Read(addr uint16) (uint8, error) {
	if addr >= 0x1000 {
		return 0, fmt.Errorf("Unknown location %d", addr)
	}

	return m.Data[addr], nil
}

func (m *Mem) ReadRange(addr uint16, len uint16) ([]uint8, error) {
	if addr+len > 0x1000 {
		return nil, fmt.Errorf("Unknown location %d", addr)
	}

	byte_arr := make([]uint8, len)
	for i := 0; i < int(len); i++ {
		byte_arr[i] = m.Data[int(addr)+i]
	}

	return byte_arr, nil
}

func (m *Mem) Write(addr uint16, v uint8) error {
	if addr >= 0x1000 {
		return fmt.Errorf("Unknown location %d", addr)
	}

	m.Data[addr] = v
	return nil
}
