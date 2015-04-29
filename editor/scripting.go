package editor

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"unicode"
)

// runs in its own goroutine?
func (e *Editor) RunInterpreter() {
}

func (e *Editor) Interpret(sexp string) ([]byte, error) {
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
		return nil, nil
	}
	sexp = strings.TrimPrefix(strings.TrimSuffix(sexp, ")"), "(")
	if len(sexp) != l-2 {
		// didn't trim off the parens
		return nil, errors.New("Malformed sexpr")
	}

	parts := strings.Split(sexp, " ")
	err := e.InterpretInternal(parts)
	if err == nil {
		return nil, nil
	}

	ans, err := runExternal(parts)
	return ans, err
}

func (e *Editor) InterpretInternal(parts []string) error {
	fn, ok := e.viewCommands[parts[0]]
	if !ok {
		return fmt.Errorf(`Interpret: "%v": function not found`, parts[0])
	}
	fn(e.views[e.currentView], 1)

	return nil
}

func runExternal(parts []string) ([]byte, error) {
	var cmd *exec.Cmd
	switch len(parts) {
	case 0:
		return nil, errors.New("runExternal cannot run empty command")
	case 1:
		cmd = exec.Command(parts[0])
	default:
		cmd = exec.Command(parts[0], parts[1:]...)

	}
	return cmd.Output()
}
