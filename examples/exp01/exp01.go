package main

import (
    "fmt"
    "time"
    "github.com/russtheaerialist/five"
)

func main() {
    board, err := five.NewBoard()
    defer board.Close()
    if err != nil {
    	panic("Unable to create board")
    }

    <- board.Ready() // Wait for the ready signal
    fmt.Println("Board Ready")

    led := board.Led(13)

    led.Strobe(0)

    fmt.Println("Waiting for 30 seconds")
    time.Sleep(time.Second * 30)  // Wait for one minute and then call stop
    fmt.Println("Stopping Strobe")

    led.Stop()

    <- board.Done() // Wait until we are "done"

    fmt.Println("Done\n\n")
}
