package chatfile

import (
	"context"
	"errors"
	"io"

	"github.com/sashabaranov/go-openai"
)

type OpenAiHistory struct {
	messages []openai.ChatCompletionMessage
}

func (h *OpenAiHistory) Append(role Role, message string) {
	var apiRole string

	switch role {
	case RoleSystem:
		apiRole = openai.ChatMessageRoleSystem
	case RoleUser:
		apiRole = openai.ChatMessageRoleUser
	case RoleAssistant:
		apiRole = openai.ChatMessageRoleAssistant
	}

	h.messages = append(h.messages, openai.ChatCompletionMessage{Role: apiRole, Content: message})
}

func (h *OpenAiHistory) PrependHistory(header OpenAiHistory) {
	h.messages = append(header.messages, h.messages...)
}

type RequestParams struct {
	Seed        *int
	Temperature float32
}

// Send a streaming request and write the response content to the provided writer in chunks as they arrive.
func Send(client *openai.Client, model ModelName, history OpenAiHistory, writer io.StringWriter, params RequestParams) (err error) {
	stream, err := client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       string(model),
			Messages:    history.messages,
			Temperature: params.Temperature,
			Seed:        params.Seed,
		},
	)

	if err != nil {
		return
	}

	defer func(stream *openai.ChatCompletionStream) {
		err = errors.Join(
			err,
			stream.Close(),
		)
	}(stream)

	for chunk, err := stream.Recv(); err == nil; chunk, err = stream.Recv() {
		if len(chunk.Choices) > 0 {
			_, err = writer.WriteString(chunk.Choices[0].Delta.Content)
		}
	}

	return
}

func NewParameters(seed *int, temperature *float32) (p RequestParams) {
	if seed != nil {
		p.Seed = seed
	}
	if temperature != nil {
		p.Temperature = *temperature
	}

	return
}
