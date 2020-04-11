package chip8

// 60 Hz. (60 opcodes in one second)

type Timer struct {
	Counter  uint8
	IsActive bool
}

func (t *Timer) Init() {
	t.Counter = 0
	t.IsActive = false
}

func (t *Timer) Set(counter_val uint8) {
	t.Counter = counter_val
	t.IsActive = t.Counter > 0
}

/// @return bool ; true if reached zero while active, false otherwise
func (t *Timer) Dec() bool {
	if t.Counter > 0 {
		t.Counter -= 1
	}

	became_zero := t.Counter == 0 && t.IsActive
	t.IsActive = t.Counter > 0

	return became_zero
}
