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
					"user": []section{
						{
							name: "user",
							props: map[string]property{
								"shell": property{
									key: "shell",
									vals: map[string][]string{
										"": []string{"/bin/bash"},
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
						"Greeting": property{
							key: "Greeting",
							vals: map[string][]string{
								"en": []string{"Hello"},
								"fr": []string{"Bonjour"},
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
					"user": []section{
						{
							name: "user",
							props: map[string]property{
								"name": property{
									key: "name",
									vals: map[string][]string{
										"": []string{"root"},
									},
								},
								"shell": property{
									key: "shell",
									vals: map[string][]string{
										"unix":  []string{"/bin/bash"},
										"win32": []string{"PowerShell.exe"},
									},
								},
							},
						},
						{
							name: "user",
							props: map[string]property{
								"name": property{
									key: "name",
									vals: map[string][]string{
										"": []string{"admin"},
									},
								},
								"shell": property{
									key: "shell",
									vals: map[string][]string{
										"unix":  []string{"/bin/bash"},
										"win32": []string{"PowerShell.exe"},
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
		input string
		want  property
	}{
		{
			input: "shell=/bin/bash",
			want: property{
				key: "shell",
				vals: map[string][]string{
					"": []string{"/bin/bash"},
				},
			},
		},
		{
			input: "Greeting[en]=Hello\nGreeting[fr]=Bonjour",
			want: property{
				key: "Greeting",
				vals: map[string][]string{
					"en": []string{"Hello"},
					"fr": []string{"Bonjour"},
				},
			},
		},
	}

	for _, test := range tests {
		p := newParser([]byte(test.input))
		p.nextToken()
		got := newProperty(p.tok.val)
		for {
			err := p.parseProperty(&got)
			if err != nil {
				t.Fatal(err)
			}
			p.nextToken()
			if p.tok.typ == tokenEOF {
				break
			}
		}
		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}

func TestParseSection(t *testing.T) {
	tests := []struct {
		input string
		want  section
	}{
		{
			input: "[user]\nname=root\nshell=/bin/bash",
			want: section{
				name: "user",
				props: map[string]property{
					"name": property{
						key: "name",
						vals: map[string][]string{
							"": []string{"root"},
						},
					},
					"shell": property{
						key: "shell",
						vals: map[string][]string{
							"": []string{"/bin/bash"},
						},
					},
				},
			},
		},
		{
			input: "[user]\n; UNIX user name\nname=root\n; Default shell\nshell=/bin/bash",
			want: section{
				name: "user",
				props: map[string]property{
					"name": property{
						key: "name",
						vals: map[string][]string{
							"": []string{"root"},
						},
					},
					"shell": property{
						key: "shell",
						vals: map[string][]string{
							"": []string{"/bin/bash"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		p := newParser([]byte(test.input))
		p.nextToken()
		got := newSection(p.tok.val)
		err := p.parseSection(&got)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}
