package chip8

func main() {
	e := MakeEmu("/home/itarato/CHECKOUT/chip8/roms/IBM Logo.ch8")
	e.Init()
	e.Run()
}
