package easybuf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// A Buffer is a direct array of bytes.
// Insertion is therefor O(n).
type Buffer struct {
	content []byte
	fname   string
}

// Load loads a buffer from a reader
func (b *Buffer) Load(from io.Reader) error {
	var buf bytes.Buffer
	n, err := io.Copy(&buf, from)

	if err != nil {
		return fmt.Errorf("buffer.Load: %d bytes read, %v", n, err)
	}
	b.content = buf.Bytes()
	return nil

}

// GetLine returns the nth line in the buffer, 0 indexed
func (b Buffer) GetLine(lineno int) (string, error) {
	return b.getLineHelper(lineno, lineno)
}

func (b Buffer) getLineHelper(lineno int, endresult int) (string, error) {
	i := bytes.IndexRune(b.content, '\n')
	if lineno == 0 {
		if i == -1 {
			i = len(b.content) - 1
		}
		return string(b.content[:i+1]), nil
	}
	if i == -1 {
		return "WAT", fmt.Errorf("Bad line request: %d", endresult)
	}
	t := b
	t.content = t.content[i+1:]
	return t.getLineHelper(lineno-1, endresult)
}

// Lines returns the number of lines in the buffer
func (b Buffer) Lines() int {
	if len(b.content) > 0 {
		return bytes.Count(b.content, []byte{'\n'}) + 1
	}
	return 0
}

func (b Buffer) Write(name string) error {
	return errors.New("easybuf.Buffer.Write: Unimplimented")
}

// WriteAt implements the io.WriterAt interface
func (b *Buffer) WriteAt(p []byte, off int64) (int, error) {
	b.content = append(
		b.content[:off],
		append(p, b.content[off:]...)...)
	return len(p), nil
}

// Delete deletes n bytes forwards from off
func (b *Buffer) Delete(n, off int) {

}

// OffsetOf takes a cursor position with origin 0,0 and returns the byte offset
// of that position in the buffer
func (b *Buffer) OffsetOf(line, column int) int {
	var lineOff int
	if line > 0 {
		lineOff = indexNth(b.content, '\n', line-1) + 1
	}

	if lineOff == -1 {
		return -1
	}

	return lineOff + column
}

func indexNth(s []byte, ch byte, n int) int {
	var seen int
	for i, c := range s {
		if c == ch {

			if seen == n {
				return i
			}
			seen++
		}
	}
	return -1
}
