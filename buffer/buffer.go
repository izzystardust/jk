// Copyright 2015 Ethan Miller. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package buffer implements buffers representing views into a piece of text.
package buffer

import "io"

type View interface {
	ByteAtOffset(n int)
	SetReadOffset(n int)
	io.Reader
}

type Buffer struct {
	Contents View
}
