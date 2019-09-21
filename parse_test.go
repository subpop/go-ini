package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  ast
	}{
		{
			input: "version=1.2.3\n\n[user]\nname=root\nshell=/bin/bash\n\n[user]\nname=admin\nshell=/bin/bash",
			want: ast{
				"": []section{
					section{
						name: "",
						props: map[string]property{
							"version": property{
								key: "version",
								val: []string{"1.2.3"},
							},
						},
					},
				},
				"user": []section{
					section{
						name: "user",
						props: map[string]property{
							"name": property{
								key: "name",
								val: []string{"root"},
							},
							"shell": property{
								key: "shell",
								val: []string{"/bin/bash"},
							},
						},
					},
					section{
						name: "user",
						props: map[string]property{
							"name": property{
								key: "name",
								val: []string{"admin"},
							},
							"shell": property{
								key: "shell",
								val: []string{"/bin/bash"},
							},
						},
					},
				},
			},
		},
		{
			input: `
version=1.2.3

[owner]
name=John Doe
organization=Acme Widgets Inc.

[database]
server=192.0.2.62
port=143
file="payroll.dat"`,
			want: ast{
				"": []section{
					section{
						name: "",
						props: map[string]property{
							"version": property{
								key: "version",
								val: []string{"1.2.3"},
							},
						},
					},
				},
				"owner": []section{
					section{
						name: "owner",
						props: map[string]property{
							"name": property{
								key: "name",
								val: []string{"John Doe"},
							},
							"organization": property{
								key: "organization",
								val: []string{"Acme Widgets Inc."},
							},
						},
					},
				},
				"database": []section{
					section{
						name: "database",
						props: map[string]property{
							"server": property{
								key: "server",
								val: []string{"192.0.2.62"},
							},
							"port": property{
								key: "port",
								val: []string{"143"},
							},
							"file": property{
								key: "file",
								val: []string{`"payroll.dat"`},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		p := newParser([]byte(test.input))
		err := p.parse()
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(p.ast, test.want, cmp.Options{cmp.AllowUnexported(section{}, property{})}) {
			t.Fatalf("%v != %v", p.ast, test.want)
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
				val: []string{
					"/bin/bash",
				},
			},
		},
	}

	for _, test := range tests {
		p := newParser([]byte(test.input))
		p.nextToken()
		got := property{
			key: p.tok.val,
			val: []string{},
		}
		err := p.parseProperty(&got)
		if err != nil {
			t.Fatal(err)
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
						val: []string{"root"},
					},
					"shell": property{
						key: "shell",
						val: []string{"/bin/bash"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		p := newParser([]byte(test.input))
		p.nextToken()
		got := section{
			name:  p.tok.val,
			props: map[string]property{},
		}
		err := p.parseSection(&got)
		if err != nil {
			t.Fatal(err)
		}
		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}
