package main

type Role string

const (
	RoleSystem    Role = "SYSTEM"
	RoleUser      Role = "USER"
	RoleAssistant Role = "ASSISTANT"
)

type ChatHistory interface {
	Append(role Role, message string)
}

type ModelName string

// Context represents the mutable state built up as commands from a chatfile
// are executed. It acts as an accumulator that collects all the information needed
// to construct a complete request to an AI service.
//
// The Context starts empty and is populated through the sequential execution of commands.
// Once all commands from a chatfile have been processed, the resulting Context contains
// the complete state needed to send the request to the AI service.
type Context struct {
	History      ChatHistory
	CurrentModel ModelName
}
