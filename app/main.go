package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/alexflint/go-arg"
	"io"
	"os"
)

type WriterHistory struct {
	writer io.Writer
}

func (w *WriterHistory) Append(role Role, message string) {
	_, _ = fmt.Fprintf(w.writer, "%s: %s\n", role, message)
}

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

	dump(file)
}

func dump(file io.Reader) {
	buffer := &bytes.Buffer{}
	history := &WriterHistory{writer: buffer}

	context := &Context{History: history}

	reader := bufio.NewReader(file)
	lexer := NewLexer(reader)
	scanner := NewParseScanner(lexer)

	for scanner.Scan() {
		command := scanner.Command()
		command.Apply(context)
	}

	fmt.Printf("Model: %s\n", context.CurrentModel)
	_, _ = buffer.WriteTo(os.Stdout)
}
