package main

import (
	"bufio"
	"chatfile/test"
	"fmt"
	"io"
	"testing"
)

func TestLexer(t *testing.T) {
	test.DoTest(t, "lex", func(t *testing.T, input io.Reader, output io.Writer) {
		reader := bufio.NewReader(input)
		lexer := NewLexer(reader)

		for lexer.MoveNext() {
			_, _ = fmt.Fprintf(output, "%v\n", lexer.Current())
		}
		fmt.Fprintln(output)
		if lexer.Err() != nil {
			_, _ = fmt.Fprintf(output, "%v\n", lexer.Err())
		}
		_, _ = fmt.Fprintf(output, "%v\n", lexer.Current())
	})
}
