package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	chatfile "github.com/vorotynsky/chatfile/lib"
)

type OpenAICredentials struct {
	APIKey  string `arg:"env:OPENAI_API_KEY,--openai-api-key,required" placeholder:"KEY" help:"OpenAICredentials API key"`
	BaseUrl string `arg:"env:OPENAI_BASE_URL,--openai-url" placeholder:"URL" help:"Custom OpenAICredentials API endpoint URL"`
	Project string `arg:"env:OPENAI_PROJECT_ID,--openai-proj" placeholder:"PROJ" help:"OpenAICredentials project identifier"`
}

func exitWithError(msg string, err error) {
	fmt.Fprintln(os.Stderr, msg, err)
	os.Exit(1)
}

func loadChatfileIntoContext(file *os.File, context *chatfile.Context) error {
	reader := bufio.NewReader(file)
	lexer := chatfile.NewLexer(reader)
	scanner := chatfile.NewParseScanner(lexer)

	for scanner.Scan() {
		command := scanner.Command()
		command.Apply(context)
	}
	return scanner.Err()
}

func createClient(credentials OpenAICredentials) openai.Client {
	options := []option.RequestOption{
		option.WithAPIKey(credentials.APIKey),
	}

	if credentials.BaseUrl != "" {
		options = append(options, option.WithBaseURL(credentials.BaseUrl))
	}

	if credentials.Project != "" {
		options = append(options, option.WithProject(credentials.Project))
	}

	client := openai.NewClient(options...)
	return client
}
