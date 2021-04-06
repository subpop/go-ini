package ini

import (
	"testing"
)

func TestNext(t *testing.T) {
	tests := []struct {
		input string
		want  struct {
			runes []rune
			pos   []int
		}
	}{
		{
			input: "abc",
			want: struct {
				runes []rune
				pos   []int
			}{
				runes: []rune{'a', 'b', 'c'},
				pos:   []int{1, 2, 3},
			},
		},
	}

	for _, test := range tests {
		l := lex(test.input)
		if len(test.input) != len(test.want.runes) {
			t.Fatalf("len(test.input) != len(test.want.runes): %v != %v", len(test.input), len(test.want.runes))
		}
		for i := 0; i < len(test.input); i++ {
			got := l.next()
			if got != test.want.runes[i] {
				t.Fatalf("%v != %v", got, test.want.runes[i])
			}
			if l.pos != test.want.pos[i] {
				t.Fatalf("%v != %v", l.pos, test.want.pos[i])
			}
		}
	}
}

func TestPrev(t *testing.T) {
	tests := []struct {
		input string
		want  struct {
			runes []rune
			pos   []int
		}
	}{
		{
			input: "abc",
			want: struct {
				runes []rune
				pos   []int
			}{
				runes: []rune{'c', 'b', 'a'},
				pos:   []int{2, 1, 0},
			},
		},
	}

	for _, test := range tests {
		l := lex(test.input)
		l.pos = len(test.input)
		l.width = 1
		if len(test.input) != len(test.want.runes) {
			t.Fatalf("len(test.input) != len(test.want.runes): %v != %v", len(test.input), len(test.want.runes))
		}
		for i := 0; i < len(test.input); i++ {
			got := l.prev()
			if got != test.want.runes[i] {
				t.Fatalf("%v != %v", got, test.want.runes[i])
			}
			if l.pos != test.want.pos[i] {
				t.Fatalf("%v != %v", l.pos, test.want.pos[i])
			}
		}
	}
}

func TestPeek(t *testing.T) {
	tests := []struct {
		input string
		want  struct {
			runes []rune
			pos   []int
		}
	}{
		{
			input: "abc",
			want: struct {
				runes []rune
				pos   []int
			}{
				runes: []rune{'b', 'b', 'b'},
				pos:   []int{1, 1, 1},
			},
		},
		{
			input: "a",
			want: struct {
				runes []rune
				pos   []int
			}{
				runes: []rune{eof},
				pos:   []int{1},
			},
		},
	}

	for _, test := range tests {
		l := lex(test.input)
		l.next()
		if len(test.input) != len(test.want.runes) {
			t.Fatalf("len(test.input) != len(test.want.runes): %v != %v", len(test.input), len(test.want.runes))
		}
		for i := 0; i < len(test.input); i++ {
			got := l.peek()
			if got != test.want.runes[i] {
				t.Fatalf("%v != %v", got, test.want.runes[i])
			}
			if l.pos != test.want.pos[i] {
				t.Fatalf("%v != %v", l.pos, test.want.pos[i])
			}
		}
	}
}

func TestRpeek(t *testing.T) {
	tests := []struct {
		input string
		want  struct {
			runes []rune
			pos   []int
		}
	}{
		{
			input: "abc",
			want: struct {
				runes []rune
				pos   []int
			}{
				runes: []rune{'c', 'c', 'c'},
				pos:   []int{3, 3, 3},
			},
		},
	}

	for _, test := range tests {
		l := lex(test.input)
		l.pos = len(test.input)
		l.width = 1
		if len(test.input) != len(test.want.runes) {
			t.Fatalf("len(test.input) != len(test.want.runes): %v != %v", len(test.input), len(test.want.runes))
		}
		for i := 0; i < len(test.input); i++ {
			got := l.rpeek()
			if got != test.want.runes[i] {
				t.Fatalf("%v != %v", got, test.want.runes[i])
			}
			if l.pos != test.want.pos[i] {
				t.Fatalf("%v != %v", l.pos, test.want.pos[i])
			}
		}
	}
}

func TestCurrent(t *testing.T) {
	tests := []struct {
		input      string
		want       string
		iterations int
	}{
		{
			input:      "abc",
			want:       "ab",
			iterations: 2,
		},
	}

	for _, test := range tests {
		l := lex(test.input)

		for i := 0; i < test.iterations; i++ {
			l.next()
		}
		got := l.current()
		if got != test.want {
			t.Fatalf("%q != %q", got, test.want)
		}
	}
}

func TestLexerNextToken(t *testing.T) {
	tests := []struct {
		description string
		input       string
		want        []token
		opts        lexerOptions
	}{
		{
			description: "simple case",
			input:       "shell=/bin/bash",
			want: []token{
				{tokenPropKey, "shell"},
				{tokenAssignment, "="},
				{tokenPropValue, "/bin/bash"},
				{tokenEOF, ""},
			},
		},
		{
			description: "section",
			input:       "[user]",
			want: []token{
				{tokenSection, "user"},
				{tokenEOF, ""},
			},
		},
		{
			description: "complete case",
			input:       "; user\n[user]\nshell=/bin/bash\ngroup=wheel",
			want: []token{
				{tokenComment, `; user`},
				{tokenSection, "user"},
				{tokenPropKey, "shell"},
				{tokenAssignment, "="},
				{tokenPropValue, "/bin/bash"},
				{tokenPropKey, "group"},
				{tokenAssignment, "="},
				{tokenPropValue, "wheel"},
				{tokenEOF, ""},
			},
		},
		{
			description: "malformed section",
			input:       "[user\nshell=/bin/bash",
			want: []token{
				{tokenError, `unexpected character: '\n', sections must be closed with a ']'`},
			},
		},
		{
			description: "empty value",
			input:       "shell=",
			want: []token{
				{tokenPropKey, "shell"},
				{tokenAssignment, "="},
				{tokenError, `unexpected character: '\x00', an assignment must be followed by one or more alphanumeric characters`},
			},
		},
		{
			description: "empty value accepted",
			input:       "shell=",
			want: []token{
				{tokenPropKey, "shell"},
				{tokenAssignment, "="},
				{tokenPropValue, ""},
				{tokenEOF, ""},
			},
			opts: lexerOptions{allowEmptyValues: true},
		},
		{
			description: "missing assignment",
			input:       "shell",
			want: []token{
				{tokenError, `unexpected character: '\x00', a property key must be followed by the assignment character ('=')`},
			},
		},
		{
			description: "whitespace multiline values",
			input:       "shell=/bin/bash\n\n /bin/zsh\ngroup=wheel",
			want: []token{
				{tokenPropKey, "shell"},
				{tokenAssignment, "="},
				{tokenPropValue, "/bin/bash\n\n /bin/zsh"},
				{tokenPropKey, "group"},
				{tokenAssignment, "="},
				{tokenPropValue, "wheel"},
				{tokenEOF, ""},
			},
			opts: lexerOptions{allowMultilineWhitespacePrefix: true},
		},
		{
			description: "escaped newline multiline values",
			input:       "shell=/bin/bash\\\n/bin/zsh",
			want: []token{
				{tokenPropKey, "shell"},
				{tokenAssignment, "="},
				{tokenPropValue, "/bin/bash\\\n/bin/zsh"},
				{tokenEOF, ""},
			},
			opts: lexerOptions{allowMultilineEscapeNewline: true},
		},
		{
			description: "map keys",
			input:       "shell[win32]=PowerShell.exe\nshell[unix]=/bin/bash\nshell[]=sh",
			want: []token{
				{tokenPropKey, "shell"},
				{tokenMapKey, "win32"},
				{tokenAssignment, "="},
				{tokenPropValue, "PowerShell.exe"},
				{tokenPropKey, "shell"},
				{tokenMapKey, "unix"},
				{tokenAssignment, "="},
				{tokenPropValue, "/bin/bash"},
				{tokenPropKey, "shell"},
				{tokenMapKey, ""},
				{tokenAssignment, "="},
				{tokenPropValue, "sh"},
				{tokenEOF, ""},
			},
		},
		{
			description: "number sign comments",
			input:       "# this is a comment",
			want: []token{
				{tokenComment, "# this is a comment"},
				{tokenEOF, ""},
			},
			opts: lexerOptions{allowNumberSignComments: true},
		},
		{
			description: "number sign comment causes error",
			input:       "# this is a comment",
			want: []token{
				{tokenError, "unexpected character: '#', comments cannot begin with '#'; consider enabling Options.AllowNumberSignComments"},
			},
		},
		{
			description: "invalid line start",
			input:       "% this is an invalid line",
			want: []token{
				{tokenError, "unexpected character: '%', lines can only begin with '[', ';', or alphanumeric characters"},
			},
		},
		{
			description: "unclosed map key",
			input:       "shell[win32",
			want: []token{
				{tokenPropKey, "shell"},
				{tokenError, "unexpected character: '\\x00', subkeys must be closed with a ']'"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			l := lex(test.input)
			l.opts = test.opts
			for i := 0; ; i++ {
				got := l.nextToken()

				if got != test.want[i] {
					t.Fatalf("nextToken() = %v, want %v", got, test.want[i])
				}
				if got.typ == tokenEOF || got.typ == tokenError {
					break
				}
			}
		})
	}
}
