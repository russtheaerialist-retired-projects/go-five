package five

import (
    "fmt"
    "time"
	"github.com/kraman/go-firmata"
)

type Led interface {
	init()
	IsOn() bool

	On()
	Off()

	Strobe(time.Duration)
	Toggle()

	Stop()
}

type led struct {
	pin byte
	board extendedBoard
	firmata *firmata.FirmataClient

	pinMode firmata.PinMode
	value bool
	stop chan bool
}

func (this *led) init() {
	this.board.Mount(this.pin, this)
	this.firmata.SetPinMode(this.pin, firmata.Output)
	this.Off()
}

func (this *led) IsOn() bool {
	return this.value
}

func (this *led) Off() {
	this.firmata.DigitalWrite(uint(this.pin), false)
	this.value = false
}

func (this *led) On() {
	this.firmata.DigitalWrite(uint(this.pin), true)
	this.value = true
}

func (this *led) Stop() {
	this.stop <- true
}

func (this *led) Toggle() {
	if this.IsOn() {
		this.Off()
	} else {
		this.On()
	}
}

func (this *led) Strobe(rate time.Duration) {
	if rate <= 0 {
		rate = 100 * time.Millisecond
	} else {
		rate = rate * time.Millisecond
	}

	go func(rate time.Duration) {
		fmt.Printf("strobing on pin %d\n", this.pin)
		for {
			select {
			case _ = <- this.stop:
				fmt.Println("received stop, ending strobe")
				return
			default:
				this.Toggle()
				time.Sleep(rate)
			}

		}
	}(rate)
}