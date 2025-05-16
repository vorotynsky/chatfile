package main

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

type TokenType string

const (
	EOF     TokenType = "<EOF>"
	UNKNOWN TokenType = "<UNKNOWN>"
	FROM    TokenType = "FROM"
	MODEL   TokenType = "MODEL"
	SYSTEM  TokenType = "SYSTEM"
	ASK     TokenType = "ASK"
	ANSWER  TokenType = "ANSWER"
	PROMPT  TokenType = "PROMPT"
)

const TabSize = 4

var (
	ErrUnknownToken      = errors.New("lexer: unknown token")
	ErrUnexpectedEOF     = errors.New("lexer: unexepected eof")
	ErrExpectedPrompt    = errors.New("lexer: prompt must be on the same line as SYSTEM/ASK/ANSWER")
	ErrExpectedModelName = errors.New("lexer: model name must be on the same line as FROM")
)

// Token represents a single lexical unit extracted during the lexical analysis process.
// It contains metadata about the token's type, content, and its location (Line and Column) in the source text.
type Token struct {
	Type    TokenType
	Content string
	Line    int
	Column  int
}

// Lexer transforms raw input text into a sequence of tokens, providing an iterator-like interface for processing input.
//
// The Lexer is designed to be memory-efficient and support processing of large inputs
// without loading the entire content into memory.
//
// The lexer is designed to be used in a forward-only manner, where tokens are consumed sequentially.
type Lexer interface {
	// MoveNext advances to the next token in the input stream.
	//
	// Prepares the next token for retrieval via Current method.
	// Must be called before first Current invocation.
	//
	// Returns true if a token is available, false when tokenization is complete.
	MoveNext() bool

	// Current returns the most recently acquired token.
	// MoveNext must be called firstly.
	Current() Token

	// Err returns any error encountered during lexical analysis.
	// Returns nil if lexing completed successfully,
	Err() error
}

const (
	s_ready = iota
	s_model
	s_prompt
)

type ReaderLexer struct {
	r     *bufio.Reader
	state int
	err   error
	ln    int
	col   int
	cur   Token
}

func NewLexer(reader *bufio.Reader) Lexer {
	return &ReaderLexer{
		r:  reader,
		ln: 1, col: 1,
		err: nil,
		cur: Token{EOF, "", 1, 1},
	}
}

func (l *ReaderLexer) Current() Token {
	return l.cur
}

func (l *ReaderLexer) Err() error {
	return l.err
}

func (l *ReaderLexer) MoveNext() bool {
	if l.err != nil {
		return false
	}

	prevLine := l.ln

	err := l.skip()

	sLn, sCol := l.ln, l.col

	if err == io.EOF {
		l.cur = Token{EOF, "", sLn, sCol}
		return false
	}

	if err != nil {
		l.err = err
		return false
	}

	switch l.state {
	case s_ready:
		word, err := l.readWord()
		if err != nil {
			l.setErr(err)
			return false
		}

		command := strings.ToUpper(word)

		switch command {
		case "FROM":
			l.cur = Token{FROM, command, sLn, sCol}
			l.state = s_model
		case "SYSTEM":
			l.cur = Token{SYSTEM, command, sLn, sCol}
			l.state = s_prompt
		case "ASK":
			l.cur = Token{ASK, command, sLn, sCol}
			l.state = s_prompt
		case "ANSWER":
			l.cur = Token{ANSWER, command, sLn, sCol}
			l.state = s_prompt
		default:
			l.cur = Token{UNKNOWN, word, sLn, sCol}
			l.err = ErrUnknownToken
			return false
		}

		return true

	case s_model:
		if prevLine != l.ln {
			l.err = ErrExpectedModelName
			l.cur = Token{UNKNOWN, "", sLn, sCol}
			return false
		}

		word, err := l.readWord()
		if err != nil {
			l.setErr(err)
			return false
		}

		l.cur = Token{MODEL, word, sLn, sCol}
		l.state = s_ready
		return true

	case s_prompt:
		if prevLine != l.ln {
			l.err = ErrExpectedPrompt
			l.cur = Token{UNKNOWN, "", sLn, sCol}
			return false
		}

		firstLine, err := l.readLine()
		if err != nil {
			l.setErr(err)
			return false
		}

		if len(firstLine) < 1 {
			l.err = ErrExpectedPrompt
			l.cur = Token{UNKNOWN, "", sLn, sCol}
			return false
		}

		firstLine = strings.TrimRightFunc(firstLine, unicode.IsSpace)

		if firstLine == "|" {
			prompt, err := l.readIndentedLines()
			if err != nil {
				l.setErr(err)
				return false
			}

			l.cur = Token{PROMPT, prompt, sLn, sCol}
			l.state = s_ready
			return true
		} else {
			l.cur = Token{PROMPT, firstLine, sLn, sCol}
			l.state = s_ready
			return true
		}
	}

	return false
}

func (l *ReaderLexer) setErr(err error) {
	if err == io.EOF {
		l.err = ErrUnexpectedEOF
		l.cur = Token{EOF, "", l.ln, l.col}
	} else {
		l.err = err
	}
}

// skip all whitespaces and new lines
func (l *ReaderLexer) skip() error {
	var err error
	var r rune

	for err == nil {
		r, _, err = l.r.ReadRune()

		if err != nil {
			break
		}

		if r == '\n' {
			l.ln++
			l.col = 1
		} else if unicode.IsSpace(r) {
			l.col += 1
		} else {
			_ = l.r.UnreadRune()
			break
		}
	}

	return err
}

func (l *ReaderLexer) readWord() (string, error) {
	var word []rune
	var err error
	var r rune

	for {
		r, _, err = l.r.ReadRune()
		if err != nil {
			break
		}
		if unicode.IsSpace(r) {
			_ = l.r.UnreadRune()
			break
		}
		word = append(word, r)
		l.col += 1
	}

	return string(word), err
}

func (l *ReaderLexer) readLine() (line string, err error) {
	line, err = l.r.ReadString('\n')
	if err != nil {
		return
	}

	drop := 0
	if line[len(line)-1] == '\n' {
		drop = 1
		if len(line) > 1 && line[len(line)-2] == '\r' {
			drop = 2
		}
		_ = l.r.UnreadByte()
	}

	line = line[:len(line)-drop]
	l.col += len(line)

	return
}

func (l *ReaderLexer) readIndentedLines() (string, error) {
	var promptBuilder strings.Builder
	var isFirstLine = true

	for {
		// Check indent level
		indentLevel := 0
		for indentLevel < TabSize {
			r, s, err := l.r.ReadRune()
			if err == io.EOF {
				if promptBuilder.Len() > 0 {
					return promptBuilder.String(), nil
				}
				return "", err
			}
			if err != nil {
				return "", err
			}

			if r == '\n' {
				l.ln++
				l.col = 1
			} else if unicode.IsSpace(r) {
				l.col += s
			}

			if !unicode.IsSpace(r) {
				l.r.UnreadRune()
				break
			}

			switch {
			case r == '\t':
				indentLevel += TabSize
			case r == '\n':
				indentLevel = 0
			default:
				indentLevel++
			}
		}

		if indentLevel < TabSize {
			break
		}

		line, err := l.readLine()
		if err == io.EOF {
			if promptBuilder.Len() > 0 {
				return promptBuilder.String(), nil
			}
			return "", err
		}
		if err != nil {
			return "", err
		}

		if !isFirstLine {
			promptBuilder.WriteRune('\n')
		}
		promptBuilder.WriteString(strings.TrimRightFunc(line, unicode.IsSpace))
		isFirstLine = false
	}

	return promptBuilder.String(), nil
}
