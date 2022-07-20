package main

import (
	"github.com/pilinsin/lontan/gui"
)

//sudo sysctl -w net.core.rmem_max=2500000
func main() {
	g := gui.New("lontan", 850, 520)
	g.Run()
}
