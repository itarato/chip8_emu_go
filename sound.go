package chip8

type Sound struct {
	Timer
}

func (s *Sound) Beep() {
	// @TODO
}

func (s *Sound) Update() {
	s.Timer.Dec()

	if s.Timer.IsZero() {
		s.Beep()
	}
}
