package main

import (
	"fmt"
	"os"

	"github.com/millere/jk/buffer"
	"github.com/nsf/termbox-go"
)

func main() {

	termbox.Init()
	defer termbox.Close()

	for i, c := range os.Args[0] {
		termbox.SetCell(i, 0, c, termbox.ColorDefault, termbox.ColorDefault)
	}

	buffer, err := buffer.BufferizeFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	y := 0
	x := 0
	for _, c := range string(buffer.Contents) {
		if c == '\n' {
			y += 1
			x = 0
			continue
		}
		termbox.SetCell(x, y, c, termbox.ColorDefault, termbox.ColorDefault)
		x += 1
	}

	for {
		termbox.Flush()
		e := termbox.PollEvent()
		switch {
		case e.Key == termbox.KeyEsc:
			return
		}
	}
	fmt.Println(buffer.Filename, len(buffer.Contents))
}
