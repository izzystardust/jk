// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/millere/jk/editor"
	"github.com/millere/jk/keys"
	"github.com/nsf/termbox-go"
)

func main() {

	err := termbox.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer termbox.Close()

	buffer, err := jk.BufferizeFile(os.Args[1])
	if err != nil {
		termbox.Close()
		fmt.Println(err)
		return
	}

	xDim, yDim := termbox.Size()

	view := jk.ViewWithBuffer(buffer, jk.Normal(), 1, 1, xDim-2, yDim-2)

	for {
		view.Draw()
		termbox.Flush()
		e := termbox.PollEvent()
		k := keys.FromTermbox(e)
		err := view.Do(k)
		if err != nil {
			return
		}
	}
}
