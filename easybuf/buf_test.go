package easybuf

import "testing"

func getBuff() Buffer {
	return Buffer{
		content: []byte(`This is a line
This is line 2
This is line 3`),
	}
}

func TestGetLine(t *testing.T) {
	b := getBuff()
	tests := []string{
		"This is a line\n",
		"This is line 2\n",
		"This is line 3",
	}

	for i, expect := range tests {
		if got, err := b.GetLine(i); got != expect || err != nil {
			if err != nil {
				t.Errorf("Case %d: error?! %v", i, err)
				continue
			}
			t.Errorf("Case %d: Got %v, expected %v", i, got, expect)
		}
	}

	if got, err := b.GetLine(5); err.Error() != "Bad line request: 5" {
		t.Errorf("Wat? Got line 5: %v", got)
	}
}

func TestLines(t *testing.T) {
	b := getBuff()
	if got := b.Lines(); got != 3 {
		t.Errorf("Got %v, expected 3", got)
	}
}

func TestWriteAt(t *testing.T) {
	// TODO: Write a better test
	b := getBuff()
	b.WriteAt([]byte("hello"), 4)
	expect := "Thishello is a line\nThis is line 2\nThis is line 3"
	if string(b.content) != expect {
		t.Errorf("Got %v, expect %v", string(b.content), expect)
	}
	b.WriteAt([]byte("bye"), b.OffsetOf(1, 0))
	expect = "Thishello is a line\nbyeThis is line 2\nThis is line 3"
	if string(b.content) != expect {
		t.Errorf("Got %v, expect %v", string(b.content), expect)
	}
}

func TestIndexNth(t *testing.T) {
	content := []byte("hello, world")
	cases := []struct {
		ch     byte
		n      int
		expect int64
	}{
		{' ', 1, -1},
		{' ', 0, 6},
		{'l', 1, 3},
		{'l', 0, 2},
	}

	for i, c := range cases {
		if got := indexNth(content, c.ch, c.n); got != c.expect {
			t.Errorf("Case %d: got %v, expected %v", i, got, c.expect)
		}
	}
}

func TestOffestOf(t *testing.T) {
	b := Buffer{
		content: []byte(`This is a line
This is line 2
This is line 3`),
	}

	cases := []struct {
		line, column int
		expect       int64
	}{
		{0, 0, 0},
		{0, 1, 1},
		{1, 0, 15},
	}

	for i, c := range cases {
		if got := b.OffsetOf(c.line, c.column); got != c.expect {
			t.Errorf("Case %d: got %v, expect %v", i, got, c.expect)
		}
	}
}
