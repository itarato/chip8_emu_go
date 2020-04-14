package main

import (
	"itarato/chip8/chip8"
	"testing"
)

func TestU16Mask(t *testing.T) {
	if chip8.U16Mask(0b1000_0110_1001_1010, 0xF000) != 0b1000 {
		t.Error("U16Mask")
	}

	if chip8.U16Mask(0b1000_0110_1001_1010, 0xFF00) != 0b1000_0110 {
		t.Error("U16Mask")
	}

	if chip8.U16Mask(0b1000_0110_1001_1010, 0x0F00) != 0b0110 {
		t.Error("U16Mask")
	}

	if chip8.U16Mask(0b1000_0110_1001_1010, 0x00FF) != 0b1001_1010 {
		t.Error("U16Mask")
	}
}

func TestBitN(t *testing.T) {
	if chip8.BitN(0b1000_0110, 0) != 1 {
		t.Error("BitN")
	}
	if chip8.BitN(0b1000_0110, 1) != 0 {
		t.Error("BitN")
	}
	if chip8.BitN(0b1000_0110, 7) != 0 {
		t.Error("BitN")
	}
}
