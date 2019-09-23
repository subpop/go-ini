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
			input: "shell=/bin/bash",
			want: []token{
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "/bin/bash"},
				{typ: tokenEOF, val: ""},
			},
		},
		{
			input: "; ignore me\n[user]\nname=root\nshell=/bin/bash\n\n[user]\nname=admin\nshell=/bin/bash\n\n[group]\nname=wheel",
			want: []token{
				{typ: tokenComment, val: "; ignore me\n"},
				{typ: tokenSection, val: "user"},
				{typ: tokenKey, val: "name"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "root"},
				{typ: tokenKey, val: "shell"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "/bin/bash"},
				{typ: tokenSection, val: "user"},
				{typ: tokenKey, val: "name"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "admin"},
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
			input: "=",
			want: []token{
				{typ: tokenError, val: "invalid character: line: 1, column: 1, '='"},
			},
		},
		{
			input: "multiline=test\\\nlines",
			want: []token{
				{typ: tokenKey, val: "multiline"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "test\\\nlines"},
				{typ: tokenEOF, val: ""},
			},
		},
		{
			input: "\tkey=has spaces",
			want: []token{
				{typ: tokenKey, val: "key"},
				{typ: tokenAssignment, val: "="},
				{typ: tokenText, val: "has spaces"},
				{typ: tokenEOF, val: ""},
			},
		},
	}

	for _, test := range tests {
		l := lex(test.input)
		var i int
		for {
			tok := l.nextToken()
			if tok != test.want[i] {
				t.Fatalf("%+v != %+v", tok, test.want[i])
			}
			if tok.typ == tokenEOF || tok.typ == tokenError {
				break
			}
			i++
		}
	}
}
