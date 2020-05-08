package ini

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// unexpectedCharErr describes an invalid rune given the lexer's current state.
type unexpectedCharErr struct {
	got rune
	msg string
}

func (e unexpectedCharErr) Error() string {
	return fmt.Sprintf("unexpected character: %q, %v", e.got, e.msg)
}

type tokenType int

const (
	tokenError tokenType = iota
	tokenPropKey
	tokenMapKey
	tokenAssignment
	tokenPropValue
	tokenSection
	tokenComment
	tokenEOF
)

const (
	eof          = rune(0)
	comment      = ';'
	sectionStart = '['
	sectionEnd   = ']'
	mapKeyStart  = '['
	mapKeyEnd    = ']'
	assignment   = '='
	eol          = '\n'
	escape       = rune(92) // backslash
	space        = ' '
	tab          = '\t'
	numberSign   = '#'
)

type stateFunc func(l *lexer) stateFunc

type token struct {
	typ tokenType
	val string
}

type lexerOptions struct {
	allowMultilineEscapeNewline    bool // support escaped newlines
	allowMultilineWhitespacePrefix bool // support space-prefixed lines
	allowEmptyValues               bool // accept empty values as valid
	allowNumberSignComments        bool // treat lines beginning with the number sign (#) as a comment
}

type lexer struct {
	input  string // the string being scanned.
	start  int    // start position of current token.
	pos    int    // current position in the input.
	width  int    // width of last rune read.
	state  stateFunc
	tokens chan token
	opts   lexerOptions
}

func lex(input string) *lexer {
	l := &lexer{
		input:  input,
		state:  lexLineStart,
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

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

// prev returns the previous rune from the input and moves the position back by
// the rune width.
func (l *lexer) prev() rune {
	l.pos -= l.width

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	return r
}

// peek returns the next rune from the input without advancing the position
func (l *lexer) peek() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

// rpeek returns the previous rune from the input without moving the position
func (l *lexer) rpeek() rune {
	if l.pos <= 0 {
		l.pos = l.width
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos-l.width:])
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

// error returns an error in the form of a stateFunc.
func (l *lexer) error(err error) stateFunc {
	l.tokens <- token{
		tokenError,
		err.Error(),
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
	case r == comment:
		return lexComment
	case r == numberSign:
		if l.opts.allowNumberSignComments {
			return lexComment
		}
		return l.error(&unexpectedCharErr{r, "comments cannot begin with '#'; consider enabling Options.AllowNumberSignComments"})
	case r == sectionStart:
		return lexSection
	case unicode.IsSpace(r):
		l.ignore()
		return lexLineStart
	case unicode.IsLetter(r) || unicode.IsDigit(r):
		return lexPropKey
	default:
		return l.error(&unexpectedCharErr{r, "lines can only begin with '[', ';', or alphanumeric characters"})
	}
}

func lexComment(l *lexer) stateFunc {
	var r rune
	for {
		r = l.peek()
		if r == eol || r == eof {
			break
		}
		l.next()
	}
	l.emit(tokenComment)
	return lexLineStart
}

func lexSection(l *lexer) stateFunc {
	var r rune
	l.ignore()
	for {
		r = l.peek()
		if r == eol || r == eof {
			return l.error(&unexpectedCharErr{r, "sections must be closed with a ']'"})
		}
		if r == sectionEnd {
			break
		}
		l.next()
	}
	l.emit(tokenSection)
	l.next()
	l.ignore()
	return lexLineStart
}

func lexPropKey(l *lexer) stateFunc {
	var r rune
	for {
		r = l.peek()
		if r == eol || r == eof {
			return l.error(&unexpectedCharErr{r, "a property key must be followed by the assignment character ('=')"})
		}
		if r == assignment || r == mapKeyStart {
			break
		}
		l.next()
	}
	l.emit(tokenPropKey)
	if r == mapKeyStart {
		return lexMapKey
	}
	return lexAssignment
}

func lexMapKey(l *lexer) stateFunc {
	var r rune
	l.next()
	l.ignore()
	for {
		r = l.peek()
		if r == eol || r == eof {
			return l.error(&unexpectedCharErr{r, "subkeys must be closed with a ']'"})
		}
		if r == mapKeyEnd {
			break
		}
		l.next()
	}
	l.emit(tokenMapKey)
	l.next()
	l.ignore()
	return lexAssignment
}

func lexAssignment(l *lexer) stateFunc {
	r := l.next()
	if r != assignment {
		panic("lexer: invalid state encountered")
	}
	l.emit(tokenAssignment)
	return lexPropValue
}

func lexPropValue(l *lexer) stateFunc {
	var r rune
	for {
		r = l.peek()
		if r == eol || r == eof {
			break
		}
		l.next()
	}
	if !l.opts.allowEmptyValues && len(l.current()) == 0 {
		l.error(&unexpectedCharErr{r, "an assignment must be followed by one or more alphanumeric characters"})
	}
	if l.opts.allowMultilineWhitespacePrefix {
		l.next()
		if unicode.IsSpace(l.peek()) {
			return lexPropValue
		}
		l.prev()
	}
	if l.opts.allowMultilineEscapeNewline {
		r := l.rpeek()
		if r == escape {
			l.next()
			return lexPropValue
		}
	}
	l.emit(tokenPropValue)
	return lexLineStart
}
