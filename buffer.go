// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package jk implements buffers representing views into a piece of text.
package jk

import (
	"fmt"
	"io/ioutil"

	"github.com/nsf/termbox-go"
)

type Buffer interface {
	GetLine(lineno int) (*Line, error)
}

type Drawer interface {
	DrawAt(x int, y int, w int, h int)
	Scroll(by int)
}

type Line struct {
	prev     *Line
	next     *Line
	Contents []byte
}

type SmallFileBuffer struct {
	Filename    string
	FirstLine   *Line
	CurrentLine *Line
	LastLine    *Line
}

func BufferizeFile(filename string) (Buffer, error) {
	a := new(SmallFileBuffer)
	a.Filename = filename
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	startOfTokenIndex := 0
	currentLine := new(Line)
	a.FirstLine = currentLine
	a.CurrentLine = currentLine
	for i, c := range contents {
		if c == '\n' {
			currentLine.Contents = contents[startOfTokenIndex : i+1]
			startOfTokenIndex = i + 1
			nextLine := new(Line)
			currentLine.next = nextLine
			nextLine.prev = currentLine
			currentLine = nextLine
		}
	}
	a.LastLine = currentLine
	return a, nil
}

func ClearBox(x int, y int, w int, h int) {
	for yi := 0; yi < h; yi++ {
		for xi := 0; xi < w; xi++ {
			termbox.SetCell(x+xi, y+yi, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
	}
}

func (a *SmallFileBuffer) GetLine(x int) (*Line, error) {
	currentLine := a.FirstLine
	i := 1
	for {
		if x == i {
			return currentLine, nil
		}
		currentLine = currentLine.next
		i += 1
		if currentLine == nil {
			return nil, fmt.Errorf("line %d does not exist", x)
		}
	}
}
