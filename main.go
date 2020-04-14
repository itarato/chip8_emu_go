package main

import (
	"itarato/chip8/chip8"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	UiScale = 8
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Test Pixel",
		Bounds: pixel.R(0, 0, 64*UiScale, 32*UiScale),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	e := chip8.MakeEmu(win, UiScale, "/home/itarato/CHECKOUT/chip8/roms/Space Invaders [David Winter].ch8")
	e.Init()
	e.Run()
}

func main() {
	pixelgl.Run(run)
}
