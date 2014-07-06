package main

import "fmt"
import "github.com/russtheaerialist/five"

func main() {
    board, err := five.NewBoard()
    defer board.Close()
    if err != nil {
    	panic("Unable to create board")
    }

    <- board.Ready() // Wait for the ready signal
    fmt.Println("Board Ready")

    led := board.Led(13)

    led.Strobe()

    <- board.Done() // Wait until we are "done"

    fmt.Println("Done\n\n")
}
