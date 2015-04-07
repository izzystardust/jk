// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package editor

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/millere/jk/keys"
	"github.com/nsf/termbox-go"
)

var LogItAll *log.Logger

type Editor struct {
	views          []*View
	currentView    *View
	modes          map[string]*Mode
	editorCommands map[string]EditorFunc
	viewCommands   map[string]ModeFunc
	log            *log.Logger
	settings       map[string]int
}

func New() *Editor {
	e := new(Editor)
	e.modes = make(map[string]*Mode)

	e.buildStandardFuncs()
	e.setStandardSettings()
	e.AddLogFile("log.txt")
	LogItAll = e.log
	e.Log("Created new editor")

	return e
}

func (e *Editor) Draw() {
	if e.currentView == nil {
		return
	}
	e.currentView.Draw()
}

func (e *Editor) Do(k keys.Keypress) error {
	//e.Log("Going to do", k)
	if e.currentView == nil {
		e.Log("currentView is nil")
		return errors.New("currentView is nil")
	}
	return e.currentView.Do(k)
}

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

func (e *Editor) setStandardSettings() {
	e.settings = map[string]int{
		"tabstop": 8,
	}
}
