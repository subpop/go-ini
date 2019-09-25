package ini

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenKey
	tokenAssignment
	tokenText
	tokenSection
	tokenComment
	tokenEOF
)

const eof = rune(0)

type stateFunc func(l *lexer) stateFunc

type token struct {
	typ tokenType
	val string
}

type lexerOptions struct {
	allowMultilineEscapeNewline    bool // support escaped newlines
	allowMultilineWhitespacePrefix bool // support space-prefixed lines
	allowEmptyValues               bool // accept empty values as valid
}

type lexer struct {
	input  string // the string being scanned.
	start  int    // start position of this item.
	pos    int    // current position in the input.
	width  int    // width of last run read.
	line   int    // current line number in the input.
	col    int    // current column in the current line.
	state  stateFunc
	tokens chan token
	opts   lexerOptions
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		state:  lexLineStart,
		line:   1,
		tokens: make(chan token, 2),
	}
	return l
}

// next returns the next rune in the input and advances the position of the lexer
// ahead by the width of the rune.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	if l.input[l.pos] == '\n' {
		l.line++
		l.col = 0
	}
	l.col++

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// prev returns the previous rune from the input and moves the position back by
// the rune width.
func (l *lexer) prev() rune {
	l.pos -= l.width
	l.col--
	if l.pos < len(l.input) && l.input[l.pos] == '\n' {
		l.line--
		l.col = 0
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	return r
}

// peek returns the next rune from the input without advancing the position
func (l *lexer) peek() rune {
	r := l.next()
	l.prev()
	return r
}

// rpeek returns the previous rune from the input without moving the position
func (l *lexer) rpeek() rune {
	r := l.prev()
	l.next()
	return r
}

// current returns the value of the current token
func (l *lexer) current() string {
	return l.input[l.start:l.pos]
}

// emit emits a token of type t, resetting the start position of the lexer to
// the current position.
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.current()}
	l.start = l.pos
}

// ignore resets the start position of the lexer to the current position, but
// does not emit a token.
func (l *lexer) ignore() {
	l.start = l.pos
}

// errorf formats and returns an error in the form of a stateFunc.
func (l *lexer) errorf(format string, args ...interface{}) stateFunc {
	l.tokens <- token{
		tokenError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

// nextToken receives the next token emitted by the lexer.
func (l *lexer) nextToken() token {
	for {
		select {
		case token := <-l.tokens:
			return token
		default:
			l.state = l.state(l)
		}
	}
}

func lexLineStart(l *lexer) stateFunc {
	r := l.next()
	switch {
	case r == eof:
		l.emit(tokenEOF)
		return nil
	case r == ';':
		return lexComment
	case r == '[':
		return lexSection
	case r == '\n':
		l.ignore()
		return lexLineStart
	case unicode.IsLetter(r) || unicode.IsDigit(r):
		return lexKey
	case r == '\t':
		l.ignore()
		return lexLineStart
	default:
		return l.errorf("unexpected input: line %v: %q", l.line, r)
	}
}

func lexComment(l *lexer) stateFunc {
	for {
		r := l.peek()
		if r == '\n' || r == eof {
			break
		}
		r = l.next()
	}
	l.emit(tokenComment)
	return lexLineStart
}

func lexSection(l *lexer) stateFunc {
	l.ignore()
	for {
		r := l.peek()
		if r == '\n' || r == eof {
			return l.errorf("unexpected input: wanted ']', got %q", r)
		}
		if r == ']' {
			break
		}
		r = l.next()
	}
	l.emit(tokenSection)
	l.next()
	l.ignore()
	return lexLineStart
}

func lexKey(l *lexer) stateFunc {
	for {
		r := l.peek()
		if r == '\n' || r == eof {
			return l.errorf("unexpected input: wanted '=', got %q", r)
		}
		if r == '=' {
			break
		}
		r = l.next()
	}
	l.emit(tokenKey)
	return lexAssignment
}

func lexAssignment(l *lexer) stateFunc {
	r := l.next()
	if r != '=' {
		l.errorf("unexpected input: wanted '=', got %q", r)
	}
	l.emit(tokenAssignment)
	return lexText
}

func lexText(l *lexer) stateFunc {
	for {
		r := l.peek()
		if r == eof {
			break
		}
		if r == '\n' {
			if l.opts.allowMultilineWhitespacePrefix {
				l.next()
				if l.peek() != ' ' {
					break
				}
			}
			if l.opts.allowMultilineEscapeNewline && l.rpeek() != '\\' {
				break
			}
		}
		l.next()
	}
	if !l.opts.allowEmptyValues && len(l.current()) == 0 {
		l.errorf("invalid token: empty value")
	}
	l.emit(tokenText)
	return lexLineStart
}
