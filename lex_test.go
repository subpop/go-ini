package ini

import (
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input string
		want  []token
	}{
		{
			input: "[user]\nname=root\nshell=/bin/bash\n\n[group]\nname=wheel",
			want: []token{
				{typ: tokenSection, val: "user"},
				{typ: tokenKey, val: "name"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "root"},
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "/bin/bash"},
				{typ: tokenSection, val: "group"},
				{typ: tokenKey, val: "name"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "wheel"},
				{typ: tokenEOF, val: ""},
			},
		},
		{
			input: "shell=/bin/bash",
			want: []token{
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "/bin/bash"},
				{typ: tokenEOF, val: ""},
			},
		},
		{
			input: "=",
			want: []token{
				{typ: tokenError, val: "ini: invalid character: line 1: '='"},
			},
		},
	}

	for _, test := range tests {
		l := lex(test.input)
		var i int
		for {
			tok := l.nextToken()
			if tok != test.want[i] {
				t.Fatalf("%v != %v", tok, test.want[i])
			}
			if tok.typ == tokenEOF || tok.typ == tokenError {
				break
			}
			i++
		}
	}
}
