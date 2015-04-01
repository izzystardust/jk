// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package editor

import (
	"errors"

	"github.com/millere/jk/keys"
)

// A ModeFunc is a function that can be executed by a keypress in a mode
type ModeFunc func(v *View, count int) error

// A Mode is used to map keypresses into events in the editor
type Mode struct {
	OnEnter  func(v *View) error
	OnExit   func(v *View) error
	EventMap map[keys.Keypress]ModeFunc
}

// Normal returns a simple normal mode for testing
func Normal(e *Editor) Mode {
	m := make(map[keys.Keypress]ModeFunc)
	m[keys.Keypress{Key: 'h'}] = func(v *View, count int) error {
		v.MoveCursor(-1, 0)
		return nil
	}
	m[keys.Keypress{Key: 'n'}] = func(v *View, count int) error {
		v.MoveCursor(0, 1)
		return nil
	}
	m[keys.Keypress{Key: 'e'}] = func(v *View, count int) error {
		v.MoveCursor(0, -1)
		return nil
	}
	m[keys.Keypress{Key: 'i'}] = func(v *View, count int) error {
		v.MoveCursor(1, 0)
		return nil
	}
	m[keys.Keypress{Key: keys.Esc}] = func(v *View, count int) error {
		return errors.New("Should quit")
	}
	m[keys.Keypress{Key: 't'}] = func(v *View, count int) error {
		v.SetMode((*v.modes)["insert"])
		return nil
	}
	m[keys.Keypress{Key: 'w'}] = func(v *View, count int) error {
		return v.back.Save("")
	}
	m[keys.Keypress{Key: '<'}] = func(v *View, count int) error {
		err := v.ExecUnderCursor(e)
		if err != nil {
			LogItAll.Println(err)
		}
		return err
	}

	return Mode{
		OnEnter:  nil,
		OnExit:   nil,
		EventMap: m,
	}
}

// Insert builds insert mode :)
func Insert() Mode {
	m := make(map[keys.Keypress]ModeFunc)
	m[keys.Keypress{Key: keys.Esc}] = func(v *View, count int) error {
		v.SetMode((*v.modes)["normal"])
		return nil
	}
	m[keys.Keypress{Key: keys.Backspace}] = func(v *View, count int) error {
		v.DeleteBackwards()
		v.MoveCursor(-1, 0)
		return nil
	}

	insertable := []rune{}
	for c := rune(0x20); c <= 0x7E; c++ {
		insertable = append(insertable, c)
	}

	for _, c := range insertable {
		cc := keys.Key(c)
		m[keys.Keypress{Key: cc}] = func(v *View, count int) error {
			for i := 0; i < count; i++ {
				v.InsertChar(rune(cc))
				v.MoveCursor(1, 0)
			}
			return nil
		}
	}
	return Mode{
		OnEnter:  nil,
		OnExit:   nil,
		EventMap: m,
	}
}

// RegisterMode registers a mode for use in the editor with a name to be referred to as
func (e *Editor) RegisterMode(name string, mode Mode) {
	e.Log("Adding mode:", name)
	e.modes[name] = &mode
	e.Log("Modes:", e.modes)
}
