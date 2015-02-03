// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jk

import "github.com/nsf/termbox-go"

type View struct {
	x, y        int
	w, h        int
	CurrentLine int
	back        Buffer
}

func ViewWithBuffer(a Buffer, x, y, w, h int) View {
	return View{
		x:           x,
		y:           y,
		w:           w,
		h:           h,
		back:        a,
		CurrentLine: 1,
	}
}

func (a *View) Draw() {
	ClearBox(a.x, a.y, a.w, a.h)
	currentLine, err := a.back.GetLine(a.CurrentLine)
	if err != nil {
		// TODO: handle error better
		panic(err)
	}
	for yi := 0; yi < a.h; yi++ {
		offset := 0
		for xi, c := range string(currentLine.Contents) {
			if xi >= a.w {
				break
			}
			if c != '\t' {
				termbox.SetCell(a.x+xi+offset, a.y+yi, c, termbox.ColorDefault, termbox.ColorDefault)
			} else {
				offset += 4
			}

		}
		currentLine = currentLine.next
		if currentLine == nil {
			break
		}
	}

}
