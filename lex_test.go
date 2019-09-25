package ini

import (
	"testing"
)

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
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "/bin/bash"},
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
			input: "; user\n[user]\nshell=/bin/bash",
			want: []token{
				{typ: tokenComment, val: `; user`},
				{typ: tokenSection, val: "user"},
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "/bin/bash"},
				{typ: tokenEOF, val: ""},
			},
		},
		{
			desc:  "malformed section",
			input: "[user\nshell=/bin/bash",
			want: []token{
				{typ: tokenError, val: `unexpected input: wanted ']', got '\n'`},
			},
		},
		{
			desc:  "empty value",
			input: "shell=",
			want: []token{
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenError, val: "invalid token: empty value"},
			},
		},
		{
			desc:  "missing assignment",
			input: "shell",
			want: []token{
				{typ: tokenError, val: `unexpected input: wanted '=', got '\x00'`},
			},
		},
		{
			desc:  "whitespace multiline values",
			input: "shell=/bin/bash\n /bin/zsh",
			want: []token{
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "/bin/bash\n /bin/zsh"},
				{typ: tokenEOF, val: ""},
			},
			opts: lexerOptions{allowMultilineWhitespacePrefix: true},
		},
		{
			desc:  "escaped newline multiline values",
			input: "shell=/bin/bash\\\n/bin/zsh",
			want: []token{
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "/bin/bash\\\n/bin/zsh"},
				{typ: tokenEOF, val: ""},
			},
			opts: lexerOptions{allowMultilineEscapeNewline: true},
		},
	}

	for _, test := range tests {
		l := lex(test.input)
		l.opts = test.opts
		var i int
		for {
			got := l.nextToken()

			if got != test.want[i] {
				t.Fatalf("%v: %+v != %+v", test.desc, got, test.want[i])
			}
			if got.typ == tokenEOF || got.typ == tokenError {
				break
			}
			i++
		}
	}
}
