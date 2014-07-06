package five

import "github.com/kraman/go-firmata"

type Board interface {
	Close()
	Led(int) Led // TODO: Replace int with an led struct
	Ready() chan int
	Done() chan int
}

type extendedBoard interface {
	Board
	Log()
}

type board struct {
	firmata *firmata.FirmataClient
	ready chan int
	done chan int
}

func NewBoard() (created_board Board, reterr error) {

	retval := new(board)
	retval.ready = make(chan int, 1)
	retval.done = make(chan int, 1)

	created_board = retval

	go func () {
		f, err := Serial(retval)
		retval.firmata = f
		reterr = err
		// Do whatever initialization
		retval.ready <- 1
	}()

	return
}

func (this *board) Log() {
	
}

func (this *board) Ready() chan int {
	return this.ready
}

func (this *board) Done() chan int {
	return this.done
}

func (this *board) Close() {
	// Close the serial port
}

func (this *board) Led(pin int) Led {
	retval := led{pin, this}

	return retval
}