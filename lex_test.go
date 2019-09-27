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
		desc  string
		input string
		want  []token
		opts  lexerOptions
	}{
		{
			desc:  "simple case",
			input: "shell=/bin/bash",
			want: []token{
				{typ: tokenPropKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "/bin/bash"},
				{typ: tokenEOF, val: ""},
			},
		},
		{
			desc:  "section",
			input: "[user]",
			want: []token{
				{typ: tokenSection, val: "user"},
				{typ: tokenEOF, val: ""},
			},
		},
		{
			desc:  "complete case",
			input: "; user\n[user]\nshell=/bin/bash\ngroup=wheel",
			want: []token{
				{typ: tokenComment, val: `; user`},
				{typ: tokenSection, val: "user"},
				{typ: tokenPropKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "/bin/bash"},
				{typ: tokenPropKey, val: "group"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "wheel"},
				{typ: tokenEOF, val: ""},
			},
		},
		{
			desc:  "malformed section",
			input: "[user\nshell=/bin/bash",
			want: []token{
				{typ: tokenError, val: `unexpected character: '\n' (expected ']')`},
			},
		},
		{
			desc:  "empty value",
			input: "shell=",
			want: []token{
				{typ: tokenPropKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenError, val: "invalid token: tokenPropValue()"},
			},
		},
		{
			desc:  "missing assignment",
			input: "shell",
			want: []token{
				{typ: tokenError, val: `unexpected character: '\x00' (expected '=')`},
			},
		},
		{
			desc:  "whitespace multiline values",
			input: "shell=/bin/bash\n /bin/zsh\ngroup=wheel",
			want: []token{
				{typ: tokenPropKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "/bin/bash\n /bin/zsh"},
				{typ: tokenPropKey, val: "group"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "wheel"},
				{typ: tokenEOF, val: ""},
			},
			opts: lexerOptions{allowMultilineWhitespacePrefix: true},
		},
		{
			desc:  "escaped newline multiline values",
			input: "shell=/bin/bash\\\n/bin/zsh",
			want: []token{
				{typ: tokenPropKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "/bin/bash\\\n/bin/zsh"},
				{typ: tokenEOF, val: ""},
			},
			opts: lexerOptions{allowMultilineEscapeNewline: true},
		},
		{
			desc:  "map keys",
			input: "shell[win32]=PowerShell.exe\nshell[unix]=/bin/bash\nshell[]=sh",
			want: []token{
				{typ: tokenPropKey, val: "shell"},
				{typ: tokenMapKey, val: "win32"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "PowerShell.exe"},
				{typ: tokenPropKey, val: "shell"},
				{typ: tokenMapKey, val: "unix"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "/bin/bash"},
				{typ: tokenPropKey, val: "shell"},
				{typ: tokenMapKey, val: ""},
				{typ: tokenAssignment, val: "="},
				{typ: tokenPropValue, val: "sh"},
				{typ: tokenEOF, val: ""},
			},
		},
	}

	for _, test := range tests {
		l := lex(test.input)
		l.opts = test.opts
		for i := 0; ; i++ {
			got := l.nextToken()

			if got != test.want[i] {
				t.Fatalf("%v: %+v != %+v", test.desc, got, test.want[i])
			}
			if got.typ == tokenEOF || got.typ == tokenError {
				break
			}
		}
	}
}
