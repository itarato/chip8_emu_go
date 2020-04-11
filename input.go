package chip8

type Input struct {
	Keys [0x10]bool
}

func (i *Input) Init() {
}

func (i *Input) UpdateState() {
}

func (i *Input) IsPressed(key_idx uint8) bool {
	if key_idx >= 0x10 {
		panic("Unknown key check")
	}

	return i.Keys[key_idx]
}
