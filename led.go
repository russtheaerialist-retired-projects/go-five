package five

import "fmt"

type Led interface {
	Strobe()
}

type led struct {
	pin int
	board Board
}

func (this led) Strobe() {
	go func() {
		fmt.Printf("strobing on pin %d\n", this.pin)
		for {
			
		}
	}()
}