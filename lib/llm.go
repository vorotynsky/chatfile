package chatfile

import (
	"context"
	"errors"
	"io"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
)

type RequestParams struct {
	Seed        param.Opt[int64]
	Temperature param.Opt[float64]
}

func NewParameters(seed *int64, temperature *float64) RequestParams {
	return RequestParams{toOpt(seed), toOpt(temperature)}
}

func toOpt[T comparable](temperature *T) param.Opt[T] {
	if temperature != nil {
		return param.NewOpt(*temperature)
	}
	return param.NullOpt[T]()
}

// Send a streaming request and write the response content to the provided writer in chunks as they arrive.
func Send(client openai.Client, model ModelName, history OpenAiHistory, writer io.StringWriter, params RequestParams) (err error) {
	stream := client.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{
		Model:       string(model),
		Messages:    history.Messages(),
		Temperature: params.Temperature,
		Seed:        params.Seed,
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
