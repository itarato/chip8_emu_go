package chip8

type Sound struct {
	Timer
}

func (s *Sound) Beep() {
	// @TODO
}

func (s *Sound) Update() {
	if s.Timer.Dec() {
		s.Beep()
	}
}
