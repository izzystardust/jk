// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package buffer implements buffers representing views into a piece of text.
package buffer

import (
	"io/ioutil"

	"github.com/nsf/termbox-go"
)

type View interface {
	ByteAtOffset(n int) byte
	DrawAt(x int, y int, w int, h int)
}

func New(file string) (View, error) {
	return BufferizeFile(file)
}

type Buffer struct {
	Contents View
}

type Line struct {
	prev     *Line
	next     *Line
	Contents []byte
}

type SmallFileBuffer struct {
	Filename  string
	FirstLine *Line
	LastLine  *Line
}

func BufferizeFile(filename string) (*SmallFileBuffer, error) {
	a := new(SmallFileBuffer)
	a.Filename = filename
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	startOfTokenIndex := 0
	currentLine := new(Line)
	a.FirstLine = currentLine
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

func (a *SmallFileBuffer) ByteAtOffset(n int) byte { return 'a' }

func (b *SmallFileBuffer) DrawAt(x int, y int, w int, h int) {
	currentLine := b.FirstLine
	for yi := 0; yi < h; yi++ {
		for xi, c := range string(currentLine.Contents) {
			if xi >= w {
				break
			}
			termbox.SetCell(x+xi, y+yi, c, termbox.ColorDefault, termbox.ColorDefault)
		}
		currentLine = currentLine.next
		if currentLine == nil {
			break
		}
	}
}
