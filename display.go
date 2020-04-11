package chip8

// 64x32, monochrome.
// Sprites: 8[w]x(1..15)[h].
// Colliding pixels (overdrawn) are XOR-ed.

type Display struct {
	Pixels [2048] /*64x32*/ uint8
}

func (d *Display) Init() {
}

func (d *Display) Draw() {

}
