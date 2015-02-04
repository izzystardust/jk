// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/millere/jk/editor"
	"github.com/nsf/termbox-go"
)

func main() {

	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer termbox.Close()

	b, err := jk.BufferizeFile(os.Args[1])
	if err != nil {
		termbox.Close()
		fmt.Println(err)
		return
	}

	xDim, yDim := termbox.Size()

	view := jk.ViewWithBuffer(b, 1, 1, xDim-2, yDim-2)

	for {
		view.Draw()
		termbox.Flush()
		e := termbox.PollEvent()
		switch {
		case e.Key == termbox.KeyEsc:
			return
		case e.Key == termbox.KeyArrowDown:
			view.FirstLine += 1
		case e.Key == termbox.KeyArrowUp:
			view.FirstLine -= 1
		case e.Ch == 'h':
			view.C.X -= 1
		case e.Ch == 'j':
			view.C.Y += 1
		case e.Ch == 'k':
			view.C.Y -= 1
		case e.Ch == 'l':
			view.C.X += 1
		}
	}
}
