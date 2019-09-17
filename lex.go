package ini

import (
	"fmt"
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

type lexer struct {
	input  string // the string being scanned.
	start  int    // start position of this item.
	pos    int    // current position in the input.
	width  int    // width of last run read.
	line   int    // current line number in the input.
	col    int    // current column in the current line.
	state  stateFunc
	tokens chan token
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		state:  lexStart,
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
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

// peek returns the next rune from the input without advancing the position
func (l *lexer) peek() rune {
	r := l.next()
	l.prev()
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

func lexStart(l *lexer) stateFunc {
	r := l.next()
	switch {
	case r == ';':
		return lexComment
	case r == '[':
		return lexSection
	case r == '\n':
		l.ignore()
		return lexStart
	case r == eof:
		l.emit(tokenEOF)
		return nil
	case r == '=':
		return l.errorf("ini: invalid character: line %v: '%v'", l.line, l.current())
	default:
		return lexKey
	}
}

func lexComment(l *lexer) stateFunc {
	for {
		r := l.next()
		if newline(r) {
			break
		}
	}
	l.emit(tokenComment)
	return lexStart
}

func lexSection(l *lexer) stateFunc {
	l.ignore()
	for {
		r := l.peek()
		if r == ']' {
			break
		}
		r = l.next()
	}
	l.emit(tokenSection)
	l.next()
	l.ignore()
	return lexStart
}

func lexKey(l *lexer) stateFunc {
	for {
		r := l.peek()
		if r == '=' {
			break
		}
		r = l.next()
	}
	l.emit(tokenKey)
	return lexAssignment
}

func lexAssignment(l *lexer) stateFunc {
	if l.next() == '=' {
		l.emit(tokenAssignment)
	} else {
		l.errorf("ini: invalid character: line %v: '%v'", l.line, l.current())
	}
	return lexText
}

func lexText(l *lexer) stateFunc {
	for {
		r := l.peek()
		if newline(r) || r == eof {
			break
		}
		r = l.next()
	}
	l.emit(tokenText)
	return lexStart
}

// newline returns true if r is equal to a newline character
func newline(r rune) bool {
	return r == '\n'
}
