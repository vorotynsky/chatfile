package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		File string `arg:"positional, required, help: open a specified file as a chatfile"`
	}
	arg.MustParse(&args)

	file, err := os.Open(args.File)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
}
