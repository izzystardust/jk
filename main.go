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

	e := editor.New()
	e.RegisterMode("normal", editor.Normal())
	e.RegisterMode("insert", editor.Insert())

	err = e.AddFile(os.Args[1])
	if err != nil {
		e.Log(err)
		return
	}

	for {
		e.Draw()
		termbox.Flush()
		v := termbox.PollEvent()
		k := keys.FromTermbox(v)
		err := e.Do(k)
		if err != nil {
			return
		}
	}
}
