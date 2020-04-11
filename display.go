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

func (d *Display) Clear() {

}

func (d *Display) DrawSprite(x uint8, y uint8, n uint8) bool {
	return false
}
