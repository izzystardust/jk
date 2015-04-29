package easybuf

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func init() {
	bufLog, _ := os.Create("buflog.txt")
	log.SetOutput(bufLog)
}

// A Buffer is a direct array of bytes.
// Insertion is therefor O(n).
type Buffer struct {
	content []byte
	fname   string
}

// Load loads a buffer from a reader
func (b *Buffer) Load(from io.Reader, name string) error {
	b.fname = name
	var buf bytes.Buffer
	n, err := io.Copy(&buf, from)

	if err != nil {
		return fmt.Errorf("buffer.Load: %d bytes read, %v", n, err)
	}
	b.content = buf.Bytes()
	log.Println(buf.Bytes())
	log.Println(bytes.Count(b.content, []byte{'\n'}))
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
	if name == "" {
		name = b.fname
	}
	out, err := os.Create(name)
	if err != nil {
		panic(err.Error() + name)
	}
	total := 0
	for total != len(b.content) {
		n, err := out.Write(b.content[total:])
		if err != nil {
			return err
		}
		total += n
	}
	return nil
}

// WriteAt implements the io.WriterAt interface
func (b *Buffer) WriteAt(p []byte, off int64) (int, error) {
	b.content = append(
		b.content[:off],
		append(p, b.content[off:]...)...)
	return len(p), nil
}

// Delete deletes n bytes forwards from off
func (b *Buffer) Delete(n, off int64) {
	if off+n > int64(len(b.content)) {
		panic("This needs to be caught")
	}
	if off < 0 {
		panic("Someone did an oops. It was probably you.")
	}
	b.content = append(b.content[:off], b.content[off+n:]...)
}

// OffsetOf takes a cursor position with origin 0,0 and returns the byte offset
// of that position in the buffer
func (b *Buffer) OffsetOf(line, column int) int64 {
	var lineOff int64
	if line > 0 {
		lineOff = indexNth(b.content, '\n', line-1) + 1
	}

	if lineOff == -1 {
		return -1
	}

	return lineOff + int64(column)
}

// indexNth indexes the nth instance of ch in s, returning -1 if there is no nth instance
// This is based on the
func indexNth(s []byte, ch byte, n int) int64 {
	var seen int
	for i, c := range s {
		if c == ch {

			if seen == n {
				return int64(i)
			}
			seen++
		}
	}
	return -1
}

func (b Buffer) Len() int {
	return len(b.content)
}

func (b Buffer) Get() (string, error) {
	return "", errors.New("easybuf.Buffer.Get(): Unimplemented")
}

func (b Buffer) FromTo(off1, off2 int64) (string, error) {
	return string(b.content[off1 : off2+1]), nil
}
