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
	parent     *Editor
	buffer     *subview
	tag        *subview
	statusArea *window.Area
	mode       *Mode
	modeName   string
	modes      *map[string]*Mode
	target     *subview
}

type subview struct {
	area      *window.Area // the area the buffer is rendered to
	C         Cursor       // the position of the cursor
	Point     *Cursor      // the position of the point, which when defined sets the selection
	back      WriteBuffer  // the backing buffer
	firstLine int          // the first line of the buffer to be displayed, for scrolling
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
	v := View{
		parent: e,
		buffer: &subview{
			area: bufarea,
			C:    Cursor{0, 0},
			back: a,
		},
		tag: &subview{
			area: tagarea,
			C:    Cursor{0, 0},
			back: tagbuf.New(),
		},

		mode:       mode,
		modeName:   m,
		modes:      &e.modes,
		statusArea: statusarea,
	}
	v.target = v.buffer
	return v, nil
}

// Draw draws a View into its area
func (v *View) Draw() {
	v.drawTag()
	v.drawBuffer()
	v.drawStatusBar()
}

func (v *View) drawTag() {
	line, _ := v.tag.back.Get()
	_, w := v.tag.area.Size()
	v.tag.area.WriteLine(line, 0, 0, w, termbox.ColorBlack, termbox.ColorWhite)
	if v.tag == v.target {
		v.tag.area.SetCursor(v.tag.C.Column, 0)
	}
}

func (v *View) drawBuffer() {
	tabStop := 4
	var tabsAtCursor int
	v.buffer.area.Clear()

	_, h := v.buffer.area.Size()
	for l := 0; l < h; l++ {
		tabs := 0
		line, err := v.buffer.back.GetLine(l + v.buffer.firstLine)
		if l+v.buffer.firstLine == v.buffer.C.Line && v.buffer.C.Column-1 >= 0 {
			tabsAtCursor = strings.Count(line[:v.buffer.C.Column-1], "\t")
		}
		if err != nil {
			break
		}
		for i, c := range line {
			if c == '\t' {
				tabs++
			}
			v.buffer.area.SetCell(i+tabStop*tabs, l, c, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	if v.buffer == v.target {
		v.buffer.area.SetCursor(v.buffer.C.Column+4*tabsAtCursor, v.buffer.C.Line-v.buffer.firstLine)
	}
}

func (v *View) drawStatusBar() {
	v.statusArea.Clear()
	_, w := v.statusArea.Size()
	modeline := fmt.Sprintf("%s [buffer %d/%d]",
		v.modeName,
		v.parent.currentView+1,
		len(v.parent.views),
	)
	v.statusArea.WriteLine(modeline, 0, 0, w, termbox.ColorBlack, termbox.ColorWhite)
}

// SetCursor sets the cursor to absolute coordinates in the file
func (v *View) SetCursor(row, column int) {
	// if buffer is empty, cursor can only be at 0,0
	if v.target.back.Len() == 0 {
		v.target.C.Column = 0
		v.target.C.Line = 0
	} else {
		total := v.target.back.Lines()
		if row >= total {
			row = total - 1
		}
		line, err := v.target.back.GetLine(row)
		if err != nil {
			row = v.target.C.Line
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
		v.target.C.Column = column
		v.target.C.Line = row
	}
	h, _ := v.target.area.Size()
	h = h - 1
	if v.target.C.Line < v.target.firstLine {
		LogItAll.Println("Scroll up. Cursor at line",
			v.target.C.Line,
			"firstLine:", v.target.firstLine,
			"h:", h)
		v.target.firstLine = v.target.C.Line
	} else if v.target.C.Line >= v.target.firstLine+h-1 {
		LogItAll.Println("Scroll Down. Cursor at line",
			v.target.C.Line,
			"firstLine:", v.target.firstLine,
			"h:", h)
		v.target.firstLine = v.target.C.Line - h + 1
	} else {
		LogItAll.Println("Not scrolling. Cursor at line",
			v.target.C.Line,
			"firstLine:", v.target.firstLine,
			"h:", h)
	}

}

// MoveCursor moves the cursor relative to where it is now
func (v *View) MoveCursor(dc, dr int) {
	v.SetCursor(v.target.C.Line+dr, v.target.C.Column+dc)
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
	off := v.target.back.OffsetOf(v.target.C.Line, v.target.C.Column)
	//LogItAll.Println("Inserting", c, "at", v.C.Line, v.C.Column, "giving an offset of", off)
	v.target.back.WriteAt([]byte{c}, off)

}

// DeleteBackwards deletes one character backwards
func (v *View) DeleteBackwards() {
	offset := v.target.back.OffsetOf(v.target.C.Line, v.target.C.Column)
	//LogItAll.Println("Delete:", v.C.Line, v.C.Column, offset-2)
	if offset < 1 {
		return
	}
	v.target.back.Delete(1, offset-1)
}

func (v *View) resultUnderCursor() ([]byte, error) {
	line, err := v.target.back.GetLine(v.target.C.Line)
	if err != nil {
		return nil, err
	}

	var i, j int

	LogItAll.Println("In line:", string(line), "C:", v.target.C.Column)
	for n, c := range string(line) {
		if n < v.target.C.Column && unicode.IsSpace(c) {
			LogItAll.Println("setting i to", n)
			i = n + 1
		}
		if n >= v.target.C.Column && unicode.IsSpace(c) {
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
	return v.parent.Interpret("(" + string(line[i:j]) + ")")
}

func (v *View) ExecInsertUnderCursor() error {
	toIns, err := v.resultUnderCursor()
	if err != nil {
		return err
	}

	v.buffer.back.WriteAt(toIns,
		v.buffer.back.OffsetOf(v.buffer.C.Line, v.buffer.C.Column))
	return nil
}

func (v *View) AlternateTag() {
	if v.target == v.buffer {
		v.target = v.tag
	} else {
		v.target = v.buffer
	}
}
