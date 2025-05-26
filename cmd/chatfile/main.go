package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	chatfile "github.com/vorotynsky/chatfile/lib"
)

func main() {
	var args struct {
		File string `arg:"positional, required, help:open a specified file as a chatfile"`

		APIKey  string `arg:"env:OPENAI_API_KEY,--openai-api-key,required" placeholder:"KEY" help:"OpenAI API key"`
		BaseUrl string `arg:"env:OPENAI_BASE_URL,--openai-url" placeholder:"URL" help:"Custom OpenAI API endpoint URL"`
		Project string `arg:"env:OPENAI_PROJECT_ID,--openai-proj" placeholder:"PROJ" help:"OpenAI project identifier"`

		Temperature *float64 `arg:"--temperature" placeholder:"TEMP" help:"Temperature for the model (this option may be removed)"`
		Seed        *int64   `arg:"--seed" placeholder:"SEED" help:"Random seed for reproducible model outputs (this option may be removed)"`
	}
	arg.MustParse(&args)

	file, err := os.Open(args.File)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	history := chatfile.OpenAiHistory{}
	context := &chatfile.Context{History: &history}

	err = processFile(file, context)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error processing file:", err)
		os.Exit(1)
	}

	client := createClient(args.APIKey, args.BaseUrl, args.Project)

	parameters := chatfile.NewParameters(args.Seed, args.Temperature)
	err = chatfile.Send(client, context.CurrentModel, history, os.Stdout, parameters)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error sending request:", err)
		os.Exit(1)
	}
}

func processFile(file *os.File, context *chatfile.Context) error {
	reader := bufio.NewReader(file)
	lexer := chatfile.NewLexer(reader)
	scanner := chatfile.NewParseScanner(lexer)

	for scanner.Scan() {
		command := scanner.Command()
		command.Apply(context)
	}
	return scanner.Err()
}

func createClient(apiKey string, baseUrl string, project string) openai.Client {
	options := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}

	if baseUrl != "" {
		options = append(options, option.WithBaseURL(baseUrl))
	}

	if project != "" {
		options = append(options, option.WithProject(project))
	}

	client := openai.NewClient(options...)
	return client
}
