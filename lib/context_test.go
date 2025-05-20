package chatfile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/vorotynsky/chatfile/test"
)

type TestHistory struct {
	w io.Writer
	c *Context
}

func (h *TestHistory) Append(role Role, message string) {
	if strings.IndexByte(message, '\n') > 0 {
		_, _ = fmt.Fprintf(h.w, "[%s] %s:\n%s\n", h.c.CurrentModel, role, message)
	} else {
		_, _ = fmt.Fprintf(h.w, "[%s] %s: %s\n", h.c.CurrentModel, role, message)
	}
}

func TestContextCollecting(t *testing.T) {
	test.DoTest(t, "collected", func(t *testing.T, input io.Reader, output io.Writer) {
		history := &TestHistory{w: output}
		context := &Context{History: history}
		history.c = context

		reader := bufio.NewReader(input)
		lexer := NewLexer(reader)
		scanner := NewParseScanner(lexer)

		for scanner.Scan() {
			command := scanner.Command()
			command.Apply(context)
		}
	})
}
