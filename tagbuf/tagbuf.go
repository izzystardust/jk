package tagbuf

import (
	"bytes"

	"github.com/millere/jk/easybuf"
)

// A Buffer is a single line containing the editable tag
type Buffer struct {
	content easybuf.Buffer
}

func New() Buffer {
	b := Buffer{}
	b.content.WriteAt([]byte("save quit"), 0)
	return b
}

func (b Buffer) Get() string {
	l, _ := b.content.GetLine(0)
	return l
}

func (b *Buffer) WriteAt(bs []byte, off int64) {
	bs_ := bytes.Replace(bs, []byte{'\n'}, []byte{' '}, -1)
	b.content.WriteAt(bs_, off)
}
