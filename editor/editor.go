// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package editor implements the functionality of the jk editor
package editor

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/millere/jk/easybuf"
	"github.com/millere/jk/keys"
	"github.com/nsf/termbox-go"
)

// LogItAll logs _everything_
var LogItAll *log.Logger

// An Editor edits shit
type Editor struct {
	views          []*View
	currentView    *View
	modes          map[string]*Mode
	editorCommands map[string]EditorFunc
	viewCommands   map[string]ModeFunc
	log            *log.Logger
	shouldQuit     bool
}

// New creates and initializes a new editor
func New() *Editor {
	e := new(Editor)
	e.modes = make(map[string]*Mode)

	e.buildStandardFuncs()
	e.AddLogFile("log.txt")
	LogItAll = e.log
	e.Log("Created new editor")

	return e
}

// Draw draws the editor to the screen
func (e *Editor) Draw() {
	if e.currentView == nil {
		return
	}
	e.currentView.Draw()
}

// Do handles events
func (e *Editor) Do(k keys.Keypress) error {
	//e.Log("Going to do", k)
	if e.currentView == nil {
		e.Log("currentView is nil")
		return errors.New("currentView is nil")
	}
	err := e.currentView.Do(k)
	if e.shouldQuit {
		return errors.New("Quitting")
	}
	return err
}

// AddFile opens the file with the given name and gives it a view
func (e *Editor) AddFile(filename string) error {
	w, h := termbox.Size()
	e.Log("Adding file:", filename)
	buffer, err := BufferizeFile(filename)
	if err != nil {
		return err
	}
	view, err := e.ViewWithBuffer(buffer, "normal", 0, 0, w, h)
	if err != nil {
		return err
	}
	e.addView(&view)

	return nil
}

// NewEmptyFile creates a view with an empty buffer
func (e *Editor) NewEmptyFile() error {
	w, h := termbox.Size()
	e.Log("Starting with empty file")
	view, err := e.ViewWithBuffer(&easybuf.Buffer{}, "normal", 0, 0, w, h)
	if err != nil {
		return err
	}
	e.addView(&view)
	return nil
}

func (e *Editor) addView(v *View) {
	e.views = append(e.views, v)
	if e.currentView == nil {
		e.currentView = v
	}
}

// Log writes to the editor's logfile
func (e *Editor) Log(things ...interface{}) {
	if e.log != nil {
		e.log.Println(things...)
	}
}

// AddLogFile sets the file the editor logs to
func (e *Editor) AddLogFile(fname string) {
	logfile, err := os.Create(fname)
	if err != nil {
		panic(err)
	}
	e.log = log.New(logfile, "jk: ", log.LstdFlags)
}

// An EditorFunc is a function that can be bound to a keypress
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

	e.viewCommands["quit"] = func(v *View, count int) error {
		e.Log("Quitting")
		e.shouldQuit = true
		return nil
	}
	e.viewCommands["save"] = func(v *View, count int) error {
		return v.buffer.back.Write("")
	}
}
