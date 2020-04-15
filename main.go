package main

import (
	"fmt"
	"itarato/chip8/chip8"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	UiScale = 8
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Chip8 Emulator",
		Bounds: pixel.R(0, 0, 64*UiScale, 32*UiScale),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var rom_name string
	if len(os.Args) >= 2 {
		rom_name = os.Args[1]
	} else {
		fmt.Printf("Missing rom. Usage: %s ROM_NAME\n", os.Args[0])
		os.Exit(1)
	}

	e := chip8.MakeEmu(win, UiScale, rom_name)
	e.Init()
	e.Run()
}

func main() {
	pixelgl.Run(run)
}
