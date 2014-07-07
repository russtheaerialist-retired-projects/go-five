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
	Toggle()
	Brightness(int)

	Strobe(time.Duration)


	Stop()
}

type led struct {
	pin byte
	board extendedBoard
	firmata *firmata.FirmataClient

	pinMode firmata.PinMode
	value int
	stop chan int
	ticker *time.Ticker
}

func (this *led) init() {
	this.board.Mount(this.pin, this)
	this.stop = make(chan int)
	this.Off()
}

func (this *led) IsOn() bool {
	return this.value != 0
}

func (this *led) setPinMode(pinMode firmata.PinMode) {
	if pinMode != this.pinMode {
		this.firmata.SetPinMode(this.pin, pinMode)
		this.pinMode = pinMode
	}
}

func (this *led) Off() {
	this.setPinMode(firmata.Ouput)
	this.firmata.DigitalWrite(uint(this.pin), false)
	this.value = false
}

func (this *led) On() {
	this.setPinMode(firmata.Output)
	this.firmata.DigitalWrite(uint(this.pin), true)
	this.value = true
}

func (this *led) Brightness(level int) {
	this.setPinMode(firmata.PWM)

	this.firmata.analogWrite(this.pin, level)
	this.value = level
}

func (this *led) Stop() {
	if this.ticker != nil {
		this.ticker.Stop()
		this.ticker = nil
		this.stop <- 1
	}
}

func (this *led) Toggle() {
	if this.IsOn() {
		this.Off()
	} else {
		this.On()
	}
}

func (this *led) Strobe(rate time.Duration) {
	if this.ticker != nil {
		return
	}

	if rate <= 0 {
		rate = 100 * time.Millisecond
	} else {
		rate = rate * time.Millisecond
	}

	this.ticker = time.NewTicker(rate)

	go func() {
		fmt.Printf("strobing on pin %d\n", this.pin)
		for {
			select {
			case <- this.stop:
				return
			case <- this.ticker.C:
				this.Toggle()
			}
		}
	}()
}