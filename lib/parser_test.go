package chatfile

import (
	"bufio"
	"fmt"
	"io"
	"testing"

	"github.com/vorotynsky/chatfile/test"
)

func TestParser(t *testing.T) {
	test.DoTest(t, "parsed", func(t *testing.T, input io.Reader, output io.Writer) {
		reader := bufio.NewReader(input)
		lexer := NewLexer(reader)
		parser := NewParseScanner(lexer)

		for parser.Scan() {
			command := parser.Command()
			_, _ = fmt.Fprintf(output, "%s: %v\n", command.Name(), command)
		}

		if parser.Err() != nil {
			_, _ = fmt.Fprintf(output, "\n%v\n%v\n", parser.Err(), lexer.Current())
		}
	})
}
