package main

import (
	"context"
	"errors"
	"io"

	"github.com/openai/openai-go"
)

type OpenAiHistory struct {
	messages []openai.ChatCompletionMessageParamUnion
}

func (h *OpenAiHistory) Append(role Role, message string) {
	switch role {
	case RoleSystem:
		h.messages = append(h.messages, openai.SystemMessage(message))
	case RoleUser:
		h.messages = append(h.messages, openai.UserMessage(message))
	case RoleAssistant:
		h.messages = append(h.messages, openai.AssistantMessage(message))
	}
}

// Send a streaming request and write the response content to the provided writer in chunks as they arrive.
func Send(client openai.Client, model ModelName, history OpenAiHistory, writer io.StringWriter) (err error) {
	stream := client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Model:    string(model),
		Messages: history.messages,
	})

	defer stream.Close()

	for stream.Next() && err == nil {
		chunk := stream.Current()

		if len(chunk.Choices) > 0 {
			_, err = writer.WriteString(chunk.Choices[0].Delta.Content)
		}
	}

	if stream.Err() != nil {
		err = errors.Join(err, stream.Err())
	}

	return
}
