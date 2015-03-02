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

	editor := jk.New()
	err = editor.AddFile(os.Args[1])

	jk.RegisterMode("normal", jk.Normal())

	for {
		editor.Draw()
		termbox.Flush()
		e := termbox.PollEvent()
		k := keys.FromTermbox(e)
		err := editor.Do(k)
		if err != nil {
			return
		}
	}
}
