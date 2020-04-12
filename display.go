package chip8

// 64x32, monochrome.
// Sprites: 8[w]x(1..15)[h].
// Colliding pixels (overdrawn) are XOR-ed.

type Display struct {
	Pixels [2048] /*64x32*/ bool
}

func (d *Display) Init() {
}

func (d *Display) Draw() {

}

func (d *Display) Clear() {
	for i, _ := range d.Pixels {
		d.Pixels[i] = false
	}
}

// W:8 H:n
func (d *Display) DrawSprite(x uint8, y uint8, n uint8, byte_arr []byte) bool {
	has_collision := false

	for offs_y := 0; offs_y < int(n); offs_y++ {
		for bit_i := 0; bit_i < 8; bit_i++ {
			prev_pixel := d.GetPixel(x+uint8(bit_i), y+uint8(offs_y))
			curr_pixel := ByteBitN(byte_arr[offs_y], uint8(bit_i)) == 1
			if prev_pixel && curr_pixel {
				has_collision = true
			}

			d.SetPixel(x+uint8(bit_i), y+uint8(offs_y), curr_pixel)
		}
	}

	return has_collision
}

func (d *Display) SetPixel(x uint8, y uint8, val bool) {
	d.Pixels[y*64+x] = val
}

func (d *Display) GetPixel(x uint8, y uint8) bool {
	return d.Pixels[y*64+x]
}
