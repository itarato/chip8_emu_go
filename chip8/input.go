package chip8

import "github.com/faiface/pixel/pixelgl"

type Input struct {
	Keys [0x10]bool
}

func (i *Input) Init() {
	for idx, _ := range i.Keys {
		i.Keys[idx] = false
	}
}

func (i *Input) UpdateState(win *pixelgl.Window) {
	buttons := [16]pixelgl.Button{
		pixelgl.Key1, pixelgl.Key2, pixelgl.Key3, pixelgl.KeyC,
		pixelgl.Key4, pixelgl.Key5, pixelgl.Key6, pixelgl.KeyD,
		pixelgl.Key7, pixelgl.Key8, pixelgl.Key9, pixelgl.KeyE,
		pixelgl.KeyA, pixelgl.Key0, pixelgl.KeyB, pixelgl.KeyF,
	}
	for idx, k := range buttons {
		i.Keys[idx] = win.Pressed(k)
	}
}

func (i *Input) IsPressed(key_idx uint8) bool {
	if key_idx >= 0x10 {
		panic("Unknown key check")
	}

	return i.Keys[key_idx]
}
