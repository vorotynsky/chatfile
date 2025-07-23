package main

import (
	"errors"
	"fmt"
	"os"

	chatfile "github.com/vorotynsky/chatfile/lib"
)

type RunCmd struct {
	File string `arg:"positional, required, help:open a specified file as a chatfile"`

	Temperature *float32 `arg:"--temperature" placeholder:"TEMP" help:"Temperature for the model (this option may be removed)"`
	Seed        *int     `arg:"--seed" placeholder:"SEED" help:"Random seed for reproducible model outputs (this option may be removed)"`

	ModelFiles map[string]string `arg:"--load-as-model,separate" placeholder:"MODEL=CHATFILE" help:"Load a file as a model with the specified name. The file will be read and parsed as a chatfile. The model name can be used in subsequent commands (such as FROM) to refer to the loaded model (this option may be removed)"`

	OpenAICredentials
}

func (cmd RunCmd) Execute() {
	file, err := os.Open(cmd.File)
	if err != nil {
		exitWithError("Error opening file:", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	history := chatfile.OpenAiHistory{}
	context := &chatfile.Context{History: &history}

	err = loadChatfileIntoContext(file, context)
	if err != nil {
		exitWithError("Error processing file:", err)
	}

	substituteModelFiles(cmd.ModelFiles, context)

	client := createClient(cmd.OpenAICredentials)
	parameters := chatfile.NewParameters(cmd.Seed, cmd.Temperature)

	err = chatfile.Send(client, context.CurrentModel, history, os.Stdout, parameters)
	if err != nil {
		exitWithError("Error sending request:", err)
	}
}

func substituteModelFiles(modelFiles map[string]string, context *chatfile.Context) {
	processedModels := make(map[string]bool)

	history := context.History.(*chatfile.OpenAiHistory)

	for {
		currentModel := string(context.CurrentModel)

		modelFilePath, found := modelFiles[currentModel]
		if !found {
			break
		}

		if processedModels[currentModel] {
			exitWithError(fmt.Sprintf("Error: Circular reference detected in model '%s'", currentModel), nil)
		}

		processedModels[currentModel] = true

		modelFile, err := os.Open(modelFilePath)
		if err != nil {
			exitWithError("Error opening model file:", err)
		}

		parentHistory := chatfile.OpenAiHistory{}
		parentContext := &chatfile.Context{History: &parentHistory}
		err = loadChatfileIntoContext(modelFile, parentContext)
		err = errors.Join(err, modelFile.Close())

		if err != nil {
			exitWithError(fmt.Sprintf("Error processing model file %s:", modelFilePath), err)
		}

		history.PrependHistory(parentHistory)
		context.CurrentModel = parentContext.CurrentModel
	}
}
