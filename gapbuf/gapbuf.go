package gapbuf

import (
	"bytes"
	"io"
)

type GapBuf struct {
	buffer    []byte // backing array of bytes
	gapStart  int    // index where gap starts (1 past last character)
	gapEnd    int    // index where characters resume after gap
	readPoint int
}

func New(size int) *GapBuf {
	a := GapBuf{
		buffer:   make([]byte, size),
		gapStart: 0,
		gapEnd:   size,
	}
	return &a
}

func FromReader(r io.Reader, gapsize int) (*GapBuf, error) {
	buf := bytes.NewBuffer(make([]byte, gapsize))
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	return &GapBuf{
		buffer:   buf.Bytes(),
		gapStart: 0,
		gapEnd:   gapsize,
	}, nil
}

// Len gives the size of the GapBuffer without the gap
func (a *GapBuf) Len() int {
	return len(a.buffer) + a.gapStart - a.gapEnd
}

// At indexes the gapbuffer
func (a *GapBuf) At(i int) byte {
	if i >= a.gapStart {
		return a.buffer[i+a.gapEnd-a.gapStart]
	} else {
		return a.buffer[i]
	}
}

// Insert inserts rune r at index i
func (a *GapBuf) Insert(r rune, i int) {
}

func (a *GapBuf) Read(buf []byte) (n int, err error) {
	i := 0
	for j := range buf {
		if i > a.Len() {
			break
		}
		buf[j] = a.At(i + a.readPoint)
		//ch := a.buffer[i+a.readPoint]
		//if i+a.readPoint >= a.gapStart {
		//	index := i + a.readPoint + a.gapEnd - a.gapStart
		//	if i >= a.Len() {
		//		fmt.Printf("breaking\n")
		//		break
		//	}
		//	ch = a.buffer[index]
		//}
		//buf[j] = ch
		i += 1
	}
	a.readPoint += i
	return i, nil
}
