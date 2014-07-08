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
	Brightness(byte)

	Strobe(time.Duration)
	Pulse(time.Duration)
	Fade(byte, time.Duration)

	Stop()
}

type led struct {
	pin byte
	board extendedBoard
	firmata *firmata.FirmataClient

	pinMode firmata.PinMode
	value byte
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
	this.setPinMode(firmata.Output)
	this.firmata.DigitalWrite(uint(this.pin), false)
	this.value = 0
}

func (this *led) On() {
	this.setPinMode(firmata.Output)
	this.firmata.DigitalWrite(uint(this.pin), true)
	this.value = 255
}

func (this *led) Brightness(level byte) {
	this.setPinMode(firmata.PWM)

	this.firmata.AnalogWrite(uint(this.pin), level)
	this.value = level
}

func (this *led) incBrightness() {
	this.Brightness(this.value + 1)
}

func (this *led) decBrightness() {
	this.Brightness(this.value - 1)
}

func (this *led) adjustBrightness(direction bool) {
	if direction {
		this.incBrightness()
	} else {
		this.decBrightness()
	}
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

func (this *led) Pulse(rate time.Duration) {
	if this.ticker != nil {
		return
	}

	this.setPinMode(firmata.PWM)

	if rate <= 0 {
		rate = time.Second
	}

	to := rate / 255 * 2
	direction := true

	this.ticker = time.NewTicker(to)

	go func() {
		for {
			select {
			case <- this.stop:
				return

			case <- this.ticker.C:
				if this.value == 0 {
					direction = true
				}

				if this.value == 255 {
					direction = false
				}

				this.adjustBrightness(direction)
			}
		}
	}()
}

func (this *led) Fade(val byte, rate time.Duration) {
	if this.ticker != nil {
		return
	}

	this.setPinMode(firmata.PWM)

    if rate <= 0 {
    	rate = time.Second
    }

    if val < 0 {
    	val = 255
    }

	var direction bool
	if this.value <= val {
		direction = true
	} else {
		direction = false
	}

	this.ticker = time.NewTicker(rate)
	go func() {
		for {
			if (direction && this.value == 255) || (!direction && this.value == 0) || (this.value == val) {
				this.Stop()
			}

			select {

				case <- this.stop:
					return

				case <- this.ticker.C:
					this.adjustBrightness(direction)
			}
		}
	}()
}