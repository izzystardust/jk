// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jk

import "github.com/nsf/termbox-go"

// a Mode maps terminal events to
type Mode map[termbox.Event]CommandFunc

// a CommandFunc takes the cursor position in a given view and changes the view
// or its corresponding buffer accordingly.
type CommandFunc func(x, y int, view *View, argument string, count int) error
