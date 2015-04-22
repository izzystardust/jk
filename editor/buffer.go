// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package editor implements buffers representing views into a piece of text.
package editor

import (
	"io"
	"os"

	"github.com/millere/jk/easybuf"
)

type WriteBuffer interface {
	io.WriterAt
	Buffer
	Write(name string) error // Writes the file to the named string
	Delete(n, off int64)     // Deletes n bytes forwards from off
	Load(from io.Reader, name string) error
}

// A Buffer holds text - these methods enable a view to display a buffer
type Buffer interface {
	GetLine(lineno int) (string, error) // The returned line shouldn't be changed
	Lines() int                         // Returns the number of lines in the buffer
	Len() int                           // The number of bytes in the buffer
	OffsetOf(line, column int) int64
}

// BufferizeFile returns a Buffer initialized with a file
func BufferizeFile(fname string) (WriteBuffer, error) {
	b := new(easybuf.Buffer)
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	err = b.Load(f, fname)
	return b, err
}
