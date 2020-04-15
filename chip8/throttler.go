package chip8

import "time"

const OptimalTick = time.Second / 240

type Throttler struct {
	T time.Time
}

func MakeThrottler() Throttler {
	return Throttler{
		T: time.Now(),
	}
}

func (t *Throttler) Throttle() {
	diff := time.Now().Sub(t.T)
	if diff < OptimalTick {
		time.Sleep(OptimalTick - diff)
	}

	t.T = time.Now()
}
