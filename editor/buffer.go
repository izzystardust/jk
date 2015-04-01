// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package editor implements buffers representing views into a piece of text.
package editor

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nsf/termbox-go"
)

// A Buffer holds text - these methods enable a view to display a buffer
type Buffer interface {
	GetLine(lineno int) (*Line, error)
	Lines() int
	Save(name string) error
}

// A Line is a single line in a linked list of lines that compose a buffer.
// TODO: Think about using other datastructures?
type Line struct {
	prev     *Line
	next     *Line
	Contents []byte
}

func (a *Line) InsertAt(offset int, toInsert []byte) {
	// do this the naive, allocating way
	// TODO: faster? less memory intensive way?

	a.Contents = append(
		a.Contents[:offset],
		append(toInsert, a.Contents[offset:]...)...)
}

func (l *Line) DeleteNAt(n int, offset int) {
	defer func() {
		if r := recover(); r != nil {
			LogItAll.Printf("DeleteNAt(%d, %d): %v\n", n, offset, r)
		}
	}()

	l.Contents = append(
		l.Contents[:offset],
		l.Contents[offset+n:]...,
	)
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
			currentLine.Contents = make([]byte, i+1-startOfTokenIndex)
			copy(currentLine.Contents, contents[startOfTokenIndex:i+1])
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

// expensive
func (b *SmallFileBuffer) Lines() int {
	i := 0
	for l := b.FirstLine; l != b.LastLine; i++ {
		l = l.next
	}
	return i + 1
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

func (a *SmallFileBuffer) Save(name string) error {
	if name == "" {
		name = a.Filename
	}

	f, err := os.Create(name)
	if err != nil {
		return err
	}

	for l := a.FirstLine; l.next != nil; l = l.next {
		written := 0
		for written != len(l.Contents) {
			n, err := f.Write(l.Contents[written:])
			if err != nil {
				LogItAll.Println(err)
			}
			written += n
		}
	}

	return nil
}
