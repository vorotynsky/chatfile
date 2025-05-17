package main

type CommandName string

type Command interface {
	Name() CommandName
}

type FromCommand struct {
	ModelName string
}

func (c *FromCommand) Name() CommandName {
	return "FROM"
}

type Role string

const (
	RoleSystem    Role = "SYSTEM"
	RoleUser      Role = "USER"
	RoleAssistant Role = "ASSISTANT"
)

type PromptCommand struct {
	Role    Role
	Message string
}

func (c *PromptCommand) Name() CommandName {
	return "PROMPT"
}
