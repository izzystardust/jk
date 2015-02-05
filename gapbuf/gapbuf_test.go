package gapbuf

import (
	"fmt"
	"io"
	"testing"
)

type repeater byte

func (a repeater) Read(bs []byte) (n int, err error) {
	for i := range bs {
		bs[i] = byte(a) + byte(i)
	}
	return len(bs), nil
}

func TestFromReader(t *testing.T) {
	a, err := FromReader(io.LimitReader(repeater('a'), 3), 5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if a.gapStart != 0 {
		t.Errorf("Gap should start at 0, is at %v", a.gapStart)
	}
	if a.gapEnd != 5 {
		t.Errorf("Gap end should be 5, is at %v", a.gapEnd)
	}
	//fmt.Println(a)
}

func TestAt(t *testing.T) {
	a, _ := FromReader(io.LimitReader(repeater('a'), 3), 5)
	for i, c := range []byte{'a', 'b', 'c'} {
		if got := a.At(i); got != c {
			t.Errorf("Got %v, expected %v", got, c)
		}
	}
}

func TestLen(t *testing.T) {
	a, err := FromReader(io.LimitReader(repeater('a'), 3), 5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	n := a.Len()
	if n != 3 {
		t.Errorf("Expected length of 3, got %v", n)
	}

}

func TestRead(t *testing.T) {
	a, err := FromReader(io.LimitReader(repeater('a'), 3), 5)
	result := make([]byte, a.Len())
	n, err := a.Read(result)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if n != 3 {
		t.Errorf("Read %v, expected 3", n)
	}
	fmt.Println(result)
}
