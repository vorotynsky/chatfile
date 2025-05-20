package chatfile

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrExpectedCommandToken = errors.New("parser: expected command token")
)

// ParseCommand parses a single command from the provided lexer.
// It advances the lexer and attempts to identify a command type based on the first token.
//
// The function returns a Command interface implementation specific to the parsed command type or an error if parsing fails.
// It returns io.EOF, which indicates successful completion of parsing.
//
// Note: For parsing multiple commands, consider using [ParseScanner], which provides an iterator-like interface.
func ParseCommand(lexer Lexer) (command Command, err error) {
	if !lexer.MoveNext() {
		return nil, errOr(lexer.Err(), io.EOF)
	}

	switch lexer.Current().Type {
	case FROM:
		return parseFrom(lexer)
	case ASK, ANSWER, SYSTEM:
		return parsePrompt(lexer)
	default:
		err = errOr(lexer.Err(), ErrExpectedCommandToken)
	}
	return
}

// ParseScanner provides an iterator-like interface for parsing multiple commands from a lexer.
// It wraps a Lexer and allows sequential processing of commands with error handling.
//
// This is the recommended way to parse commands from input as it correctly handles
// the end-of-input condition and provides a simple scanning loop pattern.
//
// Example usage:
//
//	parser := NewParseScanner(lexer)
//	for parser.Scan() {
//	    cmd := scanner.Command()
//	    // Process the command...
//	}
//	if err := scanner.Err(); err != nil {
//	    // Handle error...
//	}
type ParseScanner struct {
	lexer   Lexer
	command Command
	err     error
}

// NewParseScanner creates and initializes a new [ParseScanner] with the given lexer.
func NewParseScanner(lexer Lexer) *ParseScanner {
	return &ParseScanner{lexer, nil, nil}
}

// Scan advances the scanner to the next command in the input.
// It returns true if a command was successfully parsed, false if
// there are no more commands to parse or any error occurred.
func (p *ParseScanner) Scan() bool {
	p.command, p.err = ParseCommand(p.lexer)
	if errors.Is(p.err, io.EOF) {
		p.err = nil
		return false
	}
	return p.err == nil
}

// Command returns the current Command object after a successful Scan.
// If Scan() has not been called or returned false, the return value is undefined.
func (p *ParseScanner) Command() Command {
	return p.command
}

// Err returns any error that occurred during scanning.
// It returns nil if the scan was successful or if parsing completed normally (EOF).
// A non-nil error indicates that parsing terminated abnormally.
func (p *ParseScanner) Err() error {
	return p.err
}

func cmdFail(command TokenType) error {
	return fmt.Errorf("parser: failed to parse command %s", command)
}

func errOr(left error, right error) error {
	if left != nil {
		return left
	}
	return right
}

// assert validates that the current token in the lexer matches the expected token type.
// It's used for internal validation during parsing to ensure the parser is in a
// consistent state. If the token type doesn't match the expected type, it panics,
// indicating a programming error in the parser logic rather than a user input error.
func assert(lexer Lexer, expected TokenType) {
	if lexer.Current().Type != expected {
		panic("invalid parsing state")
	}
}

func parseFrom(lexer Lexer) (*FromCommand, error) {
	assert(lexer, FROM)

	if !lexer.MoveNext() {
		return nil, errOr(lexer.Err(), cmdFail(FROM))
	}

	assert(lexer, MODEL)

	return &FromCommand{lexer.Current().Content}, nil
}

func parsePrompt(lexer Lexer) (*PromptCommand, error) {
	var role Role

	startingToken := lexer.Current()
	switch startingToken.Type {
	case ASK:
		role = RoleUser
	case ANSWER:
		role = RoleAssistant
	case SYSTEM:
		role = RoleSystem
	default:
		panic("invalid parsing state")
	}

	if !lexer.MoveNext() {
		return nil, errOr(lexer.Err(), cmdFail(startingToken.Type))
	}

	assert(lexer, PROMPT)

	return &PromptCommand{role, lexer.Current().Content}, nil
}
