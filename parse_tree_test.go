package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseTreeAdd(t *testing.T) {
	tests := []struct {
		sections []section
		want     parseTree
	}{
		{
			sections: []section{},
			want: parseTree{
				"": []section{
					{
						name:  "",
						props: map[string]property{},
					},
				},
			},
		},
		{
			sections: []section{
				{
					name: "user",
					props: map[string]property{
						"shell": {
							key: "shell",
							vals: map[string][]string{
								"": []string{"/bin/bash"},
							},
						},
					},
				},
				{
					name: "user",
					props: map[string]property{
						"shell": {
							key: "shell",
							vals: map[string][]string{
								"": []string{"/bin/zsh"},
							},
						},
					},
				},
			},
			want: parseTree{
				"": []section{
					{name: "", props: map[string]property{}},
				},
				"user": []section{
					{
						name: "user",
						props: map[string]property{
							"shell": {
								key: "shell",
								vals: map[string][]string{
									"": []string{"/bin/bash"},
								},
							},
						},
					},
					{
						name: "user",
						props: map[string]property{
							"shell": {
								key: "shell",
								vals: map[string][]string{
									"": []string{"/bin/zsh"},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		got := newParseTree()

		for _, s := range test.sections {
			got.add(s)
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}

func TestParseTreeGet(t *testing.T) {
	tree := parseTree{
		"user": []section{
			section{
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
	}
	tests := []struct {
		input       string
		want        []section
		shouldError bool
		wantError   error
	}{
		{
			input: "user",
			want: []section{
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
		{
			input:       "group",
			want:        []section{},
			shouldError: true,
			wantError:   &missingSectionErr{"group"},
		},
	}

	for _, test := range tests {
		got, err := tree.get(test.input)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmp.Options{cmp.AllowUnexported(missingSectionErr{})}) {
				t.Errorf("%v != %v", err, test.wantError)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
				t.Errorf("%v != %v", got, test.want)
			}
		}
	}
}

func TestSectionAdd(t *testing.T) {
	tests := []struct {
		name  string
		props []property
		want  section
	}{
		{
			name: "user",
			props: []property{
				{
					key: "name",
					vals: map[string][]string{
						"": []string{"root"},
					},
				},
				{
					key: "shell",
					vals: map[string][]string{
						"": []string{"/bin/bash"},
					},
				},
				{
					key: "uid",
					vals: map[string][]string{
						"": []string{"1000", "1001"},
					},
				},
			},
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
					"uid": property{
						key: "uid",
						vals: map[string][]string{
							"": []string{"1000", "1001"},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		got := newSection(test.name)

		for _, p := range test.props {
			got.add(p)
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}

func TestSectionGet(t *testing.T) {
	sec := section{
		name: "user",
		props: map[string]property{
			"shell": property{
				key: "shell",
				vals: map[string][]string{
					"":      []string{"/bin/bash"},
					"win32": []string{"PowerShell.exe"},
				},
			},
			"username": property{
				key: "username",
				vals: map[string][]string{
					"": []string{"root"},
				},
			},
		},
	}
	tests := []struct {
		input       string
		want        *property
		shouldError bool
		wantError   error
	}{
		{
			input: "shell",
			want: &property{
				key: "shell",
				vals: map[string][]string{
					"":      []string{"/bin/bash"},
					"win32": []string{"PowerShell.exe"},
				},
			},
		},
		{
			input:       "uid",
			want:        nil,
			shouldError: true,
			wantError:   &missingPropertyErr{"uid"},
		},
	}

	for _, test := range tests {
		got, err := sec.get(test.input)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmp.Options{cmp.AllowUnexported(missingPropertyErr{})}) {
				t.Errorf("%v != %v", err, test.wantError)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
				t.Errorf("%v != %v", got, test.want)
			}
		}
	}
}

func TestPropertyAppend(t *testing.T) {
	tests := []struct {
		key    string
		values map[string][]string
		want   property
	}{
		{
			key: "shell",
			values: map[string][]string{
				"": []string{"/bin/bash", "/bin/zsh"},
			},
			want: property{
				key: "shell",
				vals: map[string][]string{
					"": []string{"/bin/bash", "/bin/zsh"},
				},
			},
		},
		{
			key: "Greeting",
			values: map[string][]string{
				"en": []string{"Hello"},
				"fr": []string{"Bonjour"},
			},
			want: property{
				key: "Greeting",
				vals: map[string][]string{
					"":   []string{},
					"en": []string{"Hello"},
					"fr": []string{"Bonjour"},
				},
			},
		},
	}

	for _, test := range tests {
		got := newProperty(test.key)

		for k, v := range test.values {
			got.append(k, v...)
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}

func TestPropertyValues(t *testing.T) {
	prop := property{
		key: "shell",
		vals: map[string][]string{
			"":      []string{"/bin/bash"},
			"win32": []string{"PowerShell.exe"},
		},
	}

	tests := []struct {
		desc        string
		input       string
		want        []string
		shouldError bool
		wantError   error
	}{
		{
			desc:  "simple",
			input: "",
			want:  []string{"/bin/bash"},
		},
		{
			desc:  "subkey",
			input: "win32",
			want:  []string{"PowerShell.exe"},
		},
		{
			desc:        "missing subkey",
			input:       "unix",
			want:        nil,
			shouldError: true,
			wantError: &missingSubkeyErr{
				p:      prop,
				subkey: "unix",
			},
		},
	}

	for _, test := range tests {
		got, err := prop.values(test.input)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmp.Options{cmp.AllowUnexported(missingSubkeyErr{}, property{})}) {
				t.Errorf("%v != %v", err, test.wantError)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(got, test.want) {
				t.Errorf("%v != %v", got, test.want)
			}
		}
	}
}
