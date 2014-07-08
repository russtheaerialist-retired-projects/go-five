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
	Pulse(time.Duration)
	Fade(int, time.Duration)

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

func (this *led) Pulse(rate time.Duration) {
	if this.ticker != nil {
		return
	}

	this.setPinMode(firmata.PWM)

	if rate <= 0 {
		rate = time.Second
	}

	to := rate / 255 * 2
	direction := 1

	this.ticker = time.NewTicker(to)

	go func() {
		for {
			select {
			case <- this.stop:
				return

			case <- this.ticker.C:
				if this.value == 0 {
					direction = 1
				}

				if this.value == 255 {
					direction = -1
				}

				this.brightness(this.value + direction)
			}
		}
	}
}

func (this *led) Fade(value int, rate time.Duration) {
	if this.ticker != nil {
		return
	}

	this.setPinMode(firmata.PWM)
}

// Led.prototype.fade = function( val, time ) {
//   // Avoid traffic jams
//   if ( this.isRunning ) {
//     return;
//   }

//   // Reset pinMode to PWM
//   this.pinMode = this.firmata.MODES.PWM;

//   var to = ( time || 1000 ) / ( (val || 255) * 2 ),
//       direction = this.value <= val ? 1 : -1;

//   priv.set( this, {
//     isOn: true,
//     isRunning: true,
//     value: this.value
//   });

//   this.interval = setInterval(function() {
//     var valueAt = this.value;

//     if ( (direction > 0 && valueAt === 255) ||
//           (direction < 0 && valueAt === 0) ||
//             valueAt === val ) {

//       this.stop();
//     } else {
//       this.brightness( valueAt + direction );
//     }
//   }.bind(this), to);

//   return this;
// };