package chip8

// 60 Hz. (60 opcodes in one second)

type Timer struct {
	Counter int8
}

func (t *Timer) Init() {
	t.Counter = -1
}

func (t *Timer) Set(counter_val int8) {
	t.Counter = counter_val
}

func (t *Timer) Dec() {
	if t.Counter >= -1 {
		t.Counter -= 1
	}
}

func (t *Timer) IsZero() bool {
	return t.Counter == 0
}
