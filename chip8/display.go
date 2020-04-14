package chip8

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// 64x32, monochrome.
// Sprites: 8[w]x(1..15)[h].
// Colliding pixels (overdrawn) are XOR-ed.

type Display struct {
	Pixels [2048] /*64x32*/ bool
	Scale  float64
}

func MakeDisplay(scale float64) Display {
	return Display{
		Scale: scale,
	}
}

func (d *Display) Init() {
	d.Clear()
}

func (d *Display) Draw(imd *imdraw.IMDraw) {
	// Move everything to Imd.
	for y := uint8(0); y < 32; y++ {
		for x := uint8(0); x < 64; x++ {
			if d.GetPixel(x, y) {
				imd.Push(pixel.V(float64(x)*d.Scale, float64(32-y)*d.Scale))
				imd.Push(pixel.V(float64(x+1)*d.Scale, float64(32-(y+1))*d.Scale))
				imd.Rectangle(0)
			}
		}
	}
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
			curr_pixel := BitN(byte_arr[offs_y], uint8(bit_i)) == 1
			prev_pixel := d.GetPixel(x+uint8(bit_i), y+uint8(offs_y))
			if prev_pixel && curr_pixel {
				has_collision = true
			}

			d.SetPixel(x+uint8(bit_i), y+uint8(offs_y), curr_pixel != prev_pixel)
		}
	}

	return has_collision
}

func (d *Display) SetPixel(x uint8, y uint8, val bool) {
	d.Pixels[uint32(y)*64+uint32(x)] = val
}

func (d *Display) GetPixel(x uint8, y uint8) bool {
	return d.Pixels[uint32(y)*64+uint32(x)]
}
