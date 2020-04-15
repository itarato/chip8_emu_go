package chip8

import (
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type Sound struct {
	Timer
}

func MakeBeep() beep.Streamer {
	dur := 24

	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			samples[i][0] = math.Sin(float64(i))
			samples[i][1] = math.Sin(float64(i))
		}

		dur -= 1

		if dur <= 0 {
			return 0, false
		} else {
			return len(samples), true
		}
	})
}

func (s *Sound) Update() {
	if s.Timer.Dec() {
		sr := beep.SampleRate(44100)
		speaker.Init(sr, sr.N(time.Second/10))
		speaker.Play(MakeBeep())
	}
}
