// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package buffer implements buffers representing views into a piece of text.
package buffer

import (
	"io"
	"io/ioutil"
)

type View interface {
	ByteAtOffset(n int)
	SetReadOffset(n int)
	io.Reader
	io.ReaderAt
}

type Buffer struct {
	Contents View
}

type Line struct {
	prev     *Line
	next     *Line
	contents []byte
}

type SmallFileBuffer struct {
	Filename string
	contents []byte
}

func BufferizeFile(filename string) (*SmallFileBuffer, error) {
	a := new(SmallFileBuffer)
	a.Filename = filename
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	a.contents = contents
	return a, nil
}
