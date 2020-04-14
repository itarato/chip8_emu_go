package chip8

type Input struct {
	Keys [0x10]bool
}

func (i *Input) Init() {
	for idx, _ := range i.Keys {
		i.Keys[idx] = false
	}
}

func (i *Input) UpdateState() {
	// @TODO
}

func (i *Input) IsPressed(key_idx uint8) bool {
	if key_idx >= 0x10 {
		panic("Unknown key check")
	}

	return i.Keys[key_idx]
}
