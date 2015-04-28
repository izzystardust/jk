package tagbuf

import (
	"bytes"

	"github.com/millere/jk/easybuf"
)

// A Buffer is a single line containing the editable tag
type Buffer struct {
	easybuf.Buffer
}

func New() *Buffer {
	b := Buffer{}
	b.Buffer.WriteAt([]byte("save quit"), 0)
	return &b
}

func (b Buffer) Get() (string, error) {
	l, _ := b.Buffer.GetLine(0)
	return l, nil
}

func (b *Buffer) WriteAt(bs []byte, off int64) (int, error) {
	bs_ := bytes.Replace(bs, []byte{'\n'}, []byte{' '}, -1)
	return b.Buffer.WriteAt(bs_, off)
}
