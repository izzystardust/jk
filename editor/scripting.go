package jk

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// runs in its own goroutine?
func (e *Editor) RunInterpreter() {
}

func (e *Editor) Interpret(sexp string) error {
	// make neat and readable
	sexp = strings.TrimSpace(sexp)
	fixws := func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		return r
	}
	sexp = strings.Map(fixws, sexp)
	l := len(sexp)

	if l == 0 {
		return nil
	}
	sexp = strings.TrimPrefix(strings.TrimSuffix(sexp, ")"), "(")
	if len(sexp) != l-2 {
		// didn't trim off the parens
		return errors.New("Malformed sexpr")
	}

	parts := strings.Split(sexp, " ")
	e.InterpretInternal(parts)
	return nil
}

func (e *Editor) InterpretInternal(parts []string) error {
	fn, ok := e.viewCommands[parts[0]]
	if !ok {
		return fmt.Errorf(`Interpret: "%v": function not found`, parts[0])
	}
	fn(e.currentView, 1)

	return nil
}
