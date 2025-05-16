package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		File string `arg:"positional, required, help:open a specified file as a chatfile"`
	}
	arg.MustParse(&args)

	file, err := os.Open(args.File)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	dumpLexer(file)
}

func dumpLexer(file *os.File) {
	reader := bufio.NewReader(file)
	lexer := NewLexer(reader)

	out := os.Stdout

	for lexer.MoveNext() {
		_, _ = fmt.Fprintf(out, "%v\n", lexer.Current())
	}
	fmt.Fprintln(out)
	if lexer.Err() != nil {
		_, _ = fmt.Fprintf(out, "%v\n", lexer.Err())
	}
	_, _ = fmt.Fprintf(out, "%v\n", lexer.Current())
}
