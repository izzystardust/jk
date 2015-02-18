// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package keys deals with nicely representing keypresses.
package keys

import (
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

type Modifier byte

const (
	Shift Modifier = 1 << iota
	Ctrl
	Alt
)

// a Key represents a keypress event
type Keypress struct {
	Mod Modifier
	Key rune
}

// These constants represent nonprinting keys on the keyboard
const (
	F1 = utf8.MaxRune + iota
	F2
	F3
	F4
	F5
	F6
	F7
	F8
	F9
	F10
	F11
	F12
	Insert
	Delete
	Home
	End
	Up
	Down
	Left
	Right
	PgDn
	PgUp
	Backspace
	Tab
	Enter
	Esc
	Space
)

// FromTermbox creates a Key from a termbox event
func FromTermbox(e termbox.Event) Keypress {
	if e.Type != termbox.EventKey {
		return Keypress{}
	}
	var k Keypress
	if e.Mod == termbox.ModAlt {
		k.Mod |= Alt
	}

	if e.Ch != 0 {
		k.Key = e.Ch
	} else {
		switch e.Key {
		case termbox.KeyF1:
			k.Key = F1
		case termbox.KeyF2:
			k.Key = F2
		case termbox.KeyF3:
			k.Key = F3
		case termbox.KeyF4:
			k.Key = F4
		case termbox.KeyF5:
			k.Key = F5
		case termbox.KeyF6:
			k.Key = F6
		case termbox.KeyF7:
			k.Key = F7
		case termbox.KeyF8:
			k.Key = F8
		case termbox.KeyF9:
			k.Key = F9
		case termbox.KeyF10:
			k.Key = F10
		case termbox.KeyF11:
			k.Key = F11
		case termbox.KeyF12:
			k.Key = F12
		case termbox.KeyInsert:
			k.Key = Insert
		case termbox.KeyDelete:
			k.Key = Delete
		case termbox.KeyHome:
			k.Key = Home
		case termbox.KeyEnd:
			k.Key = End
		case termbox.KeyPgup:
			k.Key = PgUp
		case termbox.KeyPgdn:
			k.Key = PgDn
		case termbox.KeyArrowUp:
			k.Key = Up
		case termbox.KeyArrowDown:
			k.Key = Down
		case termbox.KeyArrowLeft:
			k.Key = Left
		case termbox.KeyArrowRight:
			k.Key = Right
		case termbox.KeyBackspace:
			k.Key = Backspace
		case termbox.KeyTab:
			k.Key = Tab
		case termbox.KeyEnter:
			k.Key = Enter
		case termbox.KeyEsc:
			k.Key = Esc
		case termbox.KeySpace:
			k.Key = Space
		case termbox.KeyBackspace2:
			k.Key = Backspace
		}
	}
	return k
}
