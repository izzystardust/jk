// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package editor

import (
	"errors"

	"github.com/millere/jk/keys"
)

type ModeFunc func(v *View, count int) error

type Mode struct {
	OnEnter  func(v *View) error
	OnExit   func(v *View) error
	EventMap map[keys.Keypress]ModeFunc
}

func Normal() Mode {
	m := make(map[keys.Keypress]ModeFunc)
	m[keys.Keypress{Key: 'h'}] = func(v *View, count int) error {
		v.MoveCursor(-1, 0)
		return nil
	}
	m[keys.Keypress{Key: 'j'}] = func(v *View, count int) error {
		v.MoveCursor(0, 1)
		return nil
	}
	m[keys.Keypress{Key: 'k'}] = func(v *View, count int) error {
		v.MoveCursor(0, -1)
		return nil
	}
	m[keys.Keypress{Key: 'l'}] = func(v *View, count int) error {
		v.MoveCursor(1, 0)
		return nil
	}
	m[keys.Keypress{Key: keys.Esc}] = func(v *View, count int) error {
		return errors.New("Should quit")
	}

	return Mode{
		OnEnter:  nil,
		OnExit:   nil,
		EventMap: m,
	}
}

func (e *Editor) RegisterMode(name string, mode Mode) {
	e.Log("Adding mode:", name)
	e.modes[name] = &mode
	var modes string
	for k, _ := range e.modes {
		modes = k + " "
	}
	e.Log("Modes:", modes)
}
