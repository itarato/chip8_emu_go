package main

import (
	"itarato/chip8/chip8"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Test Pixel",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	e := chip8.MakeEmu(win, "/home/itarato/CHECKOUT/chip8/roms/IBM Logo.ch8")
	e.Init()
	e.Run()
}

func main() {
	pixelgl.Run(run)
}
