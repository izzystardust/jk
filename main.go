package main

import (
	"os"

	"github.com/nsf/termbox-go"
)

func main() {

	termbox.Init()
	defer termbox.Close()

	for i, c := range os.Args[0] {
		termbox.SetCell(i, 0, c, termbox.ColorDefault, termbox.ColorDefault)
	}

	for {
		termbox.Flush()
		e := termbox.PollEvent()
		switch {
		case e.Key == termbox.KeyEsc:
			return
		default:
			termbox.SetCell(0, 1, e.Ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}
