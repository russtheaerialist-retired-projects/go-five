package five

import "fmt"

type Led interface {
	Strobe()
	init()
}

type led struct {
	pin int
	board extendedBoard
}

func (this *led) init() {
	this.board.Mount(this.pin, this)
}

func (this *led) Strobe() {
	go func() {
		fmt.Printf("strobing on pin %d\n", this.pin)
		for {

		}
	}()
}