// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gapbuf

import "github.com/millere/jk/line"

type GapBuf struct {
	buffer    []*line.Line // backing array of bytes
	gapStart  int          // index where gap starts (1 past last character)
	gapEnd    int          // index where characters resume after gap
	readPoint int
}

func New(size int) *GapBuf {
	a := GapBuf{
		buffer:   make([]*line.Line, size),
		gapStart: 0,
		gapEnd:   size,
	}
	return &a
}

// Len gives the size of the GapBuffer without the gap
func (a *GapBuf) Len() int {
	return len(a.buffer) + a.gapStart - a.gapEnd
}

// At indexes the gapbuffer
func (a *GapBuf) At(i int) *line.Line {
	if i >= a.gapStart {
		return a.buffer[i+a.gapEnd-a.gapStart]
	} else {
		return a.buffer[i]
	}
}

// Insert inserts rune r at index i
func (a *GapBuf) Insert(r rune, i int) {
}

func (a *GapBuf) moveGap(to int) {
}
