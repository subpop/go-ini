package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  parseTree
	}{
		{
			input: "",
			want:  newParseTree(),
		},
		{
			input: "[user]\nshell=/bin/bash",
			want: parseTree{
				global: newSection(""),
				sections: map[string][]section{
					"user": {
						{
							name: "user",
							props: map[string]property{
								"shell": {
									key: "shell",
									vals: map[string][]string{
										"": {"/bin/bash"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "Greeting[en]=Hello\nGreeting[fr]=Bonjour",
			want: parseTree{
				global: section{
					name: "",
					props: map[string]property{
						"Greeting": {
							key: "Greeting",
							vals: map[string][]string{
								"en": {"Hello"},
								"fr": {"Bonjour"},
							},
						},
					},
				},
				sections: map[string][]section{},
			},
		},
		{
			input: "[user]\nname=root\nshell[unix]=/bin/bash\nshell[win32]=PowerShell.exe\n[user]\nname=admin\nshell[unix]=/bin/bash\nshell[win32]=PowerShell.exe",
			want: parseTree{
				global: newSection(""),
				sections: map[string][]section{
					"user": {
						{
							name: "user",
							props: map[string]property{
								"name": {
									key: "name",
									vals: map[string][]string{
										"": {"root"},
									},
								},
								"shell": {
									key: "shell",
									vals: map[string][]string{
										"unix":  {"/bin/bash"},
										"win32": {"PowerShell.exe"},
									},
								},
							},
						},
						{
							name: "user",
							props: map[string]property{
								"name": {
									key: "name",
									vals: map[string][]string{
										"": {"admin"},
									},
								},
								"shell": {
									key: "shell",
									vals: map[string][]string{
										"unix":  {"/bin/bash"},
										"win32": {"PowerShell.exe"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "; this is a comment",
			want:  newParseTree(),
		},
	}

	for _, test := range tests {
		p := newParser([]byte(test.input))
		err := p.parse()
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(p.tree, test.want, cmp.Options{cmp.AllowUnexported(section{}, property{}, parseTree{})}) {
			t.Fatalf("%+v != %+v", p.tree, test.want)
		}
	}
}

func TestParseProp(t *testing.T) {
	tests := []struct {
		description string
		input       string
		want        property
		shouldError bool
		wantError   error
	}{
		{
			description: "valid property",
			input:       "Greeting[en]=Hello\nGreeting[fr]=Bonjour",
			want: property{
				key: "Greeting",
				vals: map[string][]string{
					"en": {"Hello"},
					"fr": {"Bonjour"},
				},
			},
		},
		{
			description: "unexpected token, missing property value",
			input:       "Greeting=",
			want:        property{"", map[string][]string{"": {}}},
			shouldError: true,
			wantError:   &unexpectedTokenErr{token{tokenError, `unexpected character: '\x00', an assignment must be followed by one or more alphanumeric characters`}},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var err error
			p := newParser([]byte(test.input))
			p.nextToken()

			got := newProperty(p.tok.val)
			for {
				err = p.parseProperty(&got)
				if err != nil {
					break
				}
				p.nextToken()
				if p.tok.typ == tokenEOF {
					break
				}
			}
			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmp.AllowUnexported(unexpectedTokenErr{}, token{})) {
					t.Fatalf("parseProperty(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("parseProperty(%v) returned %v, want %v", test.input, err, test.wantError)
				}
				if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{})}) {
					t.Errorf("parseProperty(%v) = %v, want %v\ndiff -want +got\n%v", test.input, got, test.want, cmp.Diff(test.want, got))
				}
			}
		})
	}
}

func TestParseSection(t *testing.T) {
	tests := []struct {
		description string
		input       string
		want        section
		shouldError bool
		wantError   error
	}{
		{
			description: "valid",
			input:       "[user]\n; UNIX user name\nname=root\n; Default shell\nshell=/bin/bash",
			want: section{
				name: "user",
				props: map[string]property{
					"name":  {"name", map[string][]string{"": {"root"}}},
					"shell": {"shell", map[string][]string{"": {"/bin/bash"}}},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			p := newParser([]byte(test.input))
			p.nextToken()
			got := newSection(p.tok.val)
			err := p.parseSection(&got)

			if test.shouldError {
				if !cmp.Equal(err, test.wantError) {
					t.Fatalf("parseSection(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("parseSection(%v) returned %v, want %v", test.input, err, test.wantError)
				}
				if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
					t.Errorf("parseSection(%v) = %v, want %v\ndiff -want +got\n%v", test.input, got, test.want, cmp.Diff(test.want, got))
				}
			}
		})
	}
}
