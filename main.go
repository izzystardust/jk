// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/millere/jk/buffer"
	"github.com/nsf/termbox-go"
)

func main() {

	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer termbox.Close()

	b, err := buffer.New(os.Args[1])
	if err != nil {
		termbox.Close()
		fmt.Println(err)
		return
	}
	xDim, yDim := termbox.Size()
	b.DrawAt(1, 1, xDim-2, yDim-2)

	for {
		termbox.Flush()
		e := termbox.PollEvent()
		switch {
		case e.Key == termbox.KeyEsc:
			return
		}
	}
}
