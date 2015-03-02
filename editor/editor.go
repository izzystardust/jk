// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jk

import (
	"fmt"
	"log"
	"os"

	"github.com/millere/jk/keys"
)

var (
	Tabstop = 4
)

type Editor struct {
	views          []*View
	currentView    *View
	modes          map[string]Mode
	editorCommands map[string]EditorFunc
	viewCommands   map[string]ModeFunc
	log            *log.Logger
}

func New() *Editor {
	e := new(Editor)

	e.buildStandardFuncs()
	e.AddLogFile("log.txt")

	return e
}

func (e *Editor) Draw() {
	if e.currentView == nil {
		return
	}
	e.currentView.Draw()
}

func (e *Editor) Do(k keys.Keypress) error {
	return e.currentView.Do(k)
}

func (e *Editor) AddFile(filename string) error {
	buffer, err := BufferizeFile(filename)
	if err != nil {
		return err
	}
	view, err := ViewWithBuffer(buffer, "normal", 1, 1, 80, 80)
	if err != nil {
		return err
	}
	e.views = append(e.views, &view)
	if e.currentView == nil {
		e.currentView = &view
	}

	return nil
}

func (e *Editor) Log(things ...interface{}) {
	if e.log != nil {
		e.log.Println(things...)
	}
}

func (e *Editor) AddLogFile(fname string) {
	logfile, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	e.log = log.New(logfile, "jk: ", log.LstdFlags)
}

type EditorFunc func(e *Editor, args ...string) error

func (e *Editor) buildStandardFuncs() {
	e.editorCommands = make(map[string]EditorFunc)
	e.viewCommands = make(map[string]ModeFunc)
	e.editorCommands["bind-key-in-mode"] = func(e *Editor, args ...string) error {
		mode := args[0]
		m, ok := e.modes[mode]
		if !ok {
			return fmt.Errorf("BindKeyInMode: no such mode %s", mode)
		}

		e.log.Println(m)
		return nil
	}
}
