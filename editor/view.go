// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package editor

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/millere/jk/keys"
	"github.com/millere/jk/tagbuf"
	"github.com/millere/window"
	"github.com/nsf/termbox-go"
)

// A View contains a buffer and knows how to draw it to an area
type View struct {
	bufarea    *window.Area
	statusArea *window.Area
	tagarea    *window.Area
	FirstLine  int         // index of the first line
	back       WriteBuffer // the backing buffer being displayed
	tag        tagbuf.Buffer
	C          Cursor
	mode       *Mode
	modeName   string
	modes      *map[string]*Mode
}

// A Cursor indicates where the cursor is
// 0, 0 is the first position in a file
type Cursor struct {
	Line, Column int
}

// ViewWithBuffer creates a new view and area with the given buffer in the editor.
func (e *Editor) ViewWithBuffer(a WriteBuffer, m string, x, y, w, h int) (View, error) {
	mode, ok := e.modes[m]
	if !ok {
		return View{}, fmt.Errorf("Mode \"%v\" does not exist", m)
	}
	tagarea := window.New(x, y, w, 1)
	bufarea := window.New(x, y+1, w, h-1)
	statusarea := window.New(x, y+h-1, w, 1)
	return View{
		back:      a,
		tag:       tagbuf.New(),
		FirstLine: 0,
		C: Cursor{
			Line:   0,
			Column: 0,
		},
		mode:       mode,
		modeName:   m,
		modes:      &e.modes,
		bufarea:    bufarea,
		statusArea: statusarea,
		tagarea:    tagarea,
	}, nil
}

// Draw draws a View into its area
func (v *View) Draw() {
	v.drawTag()
	v.drawBuffer()
	v.drawStatusBar()
}

func (v *View) drawTag() {
	line := v.tag.Get()
	_, w := v.tagarea.Size()
	v.tagarea.WriteLine(line, 0, 0, w, termbox.ColorBlack, termbox.ColorWhite)
}

func (v *View) drawBuffer() {
	tabStop := 4
	var tabsAtCursor int
	v.bufarea.Clear()

	_, h := v.bufarea.Size()
	for l := 0; l < h; l++ {
		tabs := 0
		line, err := v.back.GetLine(l + v.FirstLine)
		if l+v.FirstLine == v.C.Line && v.C.Column-1 >= 0 {
			tabsAtCursor = strings.Count(line[:v.C.Column-1], "\t")
		}
		if err != nil {
			break
		}
		for i, c := range line {
			if c == '\t' {
				tabs++
			}
			v.bufarea.SetCell(i+tabStop*tabs, l, c, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	v.bufarea.SetCursor(v.C.Column+4*tabsAtCursor, v.C.Line)
}

func (v *View) drawStatusBar() {
	v.statusArea.Clear()
	_, w := v.statusArea.Size()
	v.statusArea.WriteLine(v.modeName, 0, 0, w, termbox.ColorBlack, termbox.ColorWhite)
}

// SetCursor sets the cursor to absolute coordinates in the file
func (v *View) SetCursor(row, column int) {
	// if buffer is empty, cursor can only be at 0,0
	if v.back.Len() == 0 {
		v.C.Column = 0
		v.C.Line = 0
	} else {
		total := v.back.Lines()
		if row >= total {
			row = total - 1
		}
		line, err := v.back.GetLine(row)
		if err != nil {
			row = v.C.Line
		}
		l := len(line)
		if l > 0 && line[l-1] == '\n' {
			l -= 1
		}
		if column > l {
			column = l
		}
		if column < 0 {
			column = 0
		}
		LogItAll.Printf("Setting cursor to (%d, %d) (total: %d)", v.C.Column, v.C.Line, total)
		v.C.Column = column
		v.C.Line = row
	}
}

// MoveCursor moves the cursor relative to where it is now
func (v *View) MoveCursor(dc, dr int) {
	v.SetCursor(v.C.Line+dr, v.C.Column+dc)
}

func inBounds(x, y, w, h, ax, ay int) bool {
	return ax >= x && ax < w && ay >= y && ay < h
}

// SetMode sets a view's mode, so that it handles events per that mode
func (v *View) SetMode(m *Mode, n string) {
	if v.mode.OnExit != nil {
		v.mode.OnExit(v)
	}
	v.mode = m
	v.modeName = n
	if v.mode.OnEnter != nil {
		v.mode.OnEnter(v)
	}
}

// Do tells a view to handle a keypress according to its mode
func (v *View) Do(k keys.Keypress) error {
	f, ok := v.mode.EventMap[k]
	if ok {
		return f(v, 1)
	}
	LogItAll.Printf("No function bound to key %v", k)
	return nil

}

// InsertChar inserts the single rune r at the cursor
func (v *View) InsertChar(c byte) {
	off := v.back.OffsetOf(v.C.Line, v.C.Column)
	//LogItAll.Println("Inserting", c, "at", v.C.Line, v.C.Column, "giving an offset of", off)
	v.back.WriteAt([]byte{c}, off)
}

// DeleteBackwards deletes one character backwards
func (v *View) DeleteBackwards() {
	offset := v.back.OffsetOf(v.C.Line, v.C.Column)
	LogItAll.Println("Delete:", v.C.Line, v.C.Column, offset-2)
	v.back.Delete(1, offset-2)
}

func (v *View) ExecUnderCursor(e *Editor) error {
	line, err := v.back.GetLine(v.C.Line)
	if err != nil {
		return err
	}

	var i, j int

	LogItAll.Println("In line:", string(line), "C:", v.C.Column)
	for n, c := range string(line) {
		if n < v.C.Column && unicode.IsSpace(c) {
			LogItAll.Println("setting i to", n)
			i = n + 1
		}
		if n >= v.C.Column && unicode.IsSpace(c) {
			LogItAll.Println("setting j to", n)
			j = n
			break
		}
	}
	// if no space after...
	if j == 0 {
		j = len(line)
	}

	LogItAll.Println("Command:", string(line[i:j]), "i:", i, "j:", j)
	toIns, err := e.Interpret("(" + string(line[i:j]) + ")")
	if err != nil {
		return err
	}

	v.back.WriteAt(toIns, v.back.OffsetOf(v.C.Line, v.C.Column))
	//line.InsertAt(v.C.Column-1, toIns)
	return nil
}
