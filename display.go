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
	Imd    *imdraw.IMDraw
}

func MakeDisplay() Display {
	return Display{
		Imd: imdraw.New(nil),
	}
}

func (d *Display) Init() {
	d.Clear()
	d.Imd.Clear()
	d.Imd.Color = pixel.RGB(1, 1, 1)
}

func (d *Display) Draw() {
	// Move everything to Imd.
	d.Imd.Clear()

	for y := uint8(0); y < 32; y++ {
		for x := uint8(0); x < 64; x++ {
			pixel_on := d.GetPixel(x, y)
			if pixel_on {
				d.Imd.Push(pixel.V(float64(x)*10.0, float64(y)*10.0))
				d.Imd.Push(pixel.V(float64(x)*10.0+10.0, float64(y)*10.0+10.0))
				d.Imd.Rectangle(0)
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
