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

	dump(file)
}

func dump(file *os.File) {
	reader := bufio.NewReader(file)
	lexer := NewLexer(reader)
	scanner := NewParseScanner(lexer)

	for scanner.Scan() {
		command := scanner.Command()
		fmt.Printf("%s: %v\n", command.Name(), command)
	}

	fmt.Println()
	fmt.Println(scanner.Err())
}
