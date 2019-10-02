package ini

import (
	"reflect"
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
				global: section{
					name:  "",
					props: make(map[string]property),
				},
				sections: map[string][]section{},
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
				global: section{
					name:  "",
					props: map[string]property{},
				},
				sections: map[string][]section{
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
		},
	}

	for _, test := range tests {
		got := newParseTree()

		for _, s := range test.sections {
			got.add(s)
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{}, parseTree{})}) {
			t.Errorf("%+v != %+v", got, test.want)
		}
	}
}

func TestParseTreeGet(t *testing.T) {
	tree := parseTree{
		global: newSection(""),
		sections: map[string][]section{
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
			"root": []section{
				section{
					name: "root",
					props: map[string]property{
						"username": property{
							key: "username",
							vals: map[string][]string{
								"": []string{"root"},
							},
						},
					},
				},
			},
			"admin": []section{
				section{
					name: "admin",
					props: map[string]property{
						"username": property{
							key: "username",
							vals: map[string][]string{
								"": []string{"admin"},
							},
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
			input:       "",
			want:        []section{},
			shouldError: true,
			wantError:   &invalidKeyErr{"section name cannot be empty"},
		},
		{
			input: "*",
			want: []section{
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
				section{
					name: "root",
					props: map[string]property{
						"username": property{
							key: "username",
							vals: map[string][]string{
								"": []string{"root"},
							},
						},
					},
				},
				section{
					name: "admin",
					props: map[string]property{
						"username": property{
							key: "username",
							vals: map[string][]string{
								"": []string{"admin"},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		got, err := tree.get(test.input)

		if test.shouldError {
			if !reflect.DeepEqual(err, test.wantError) {
				t.Errorf("%v != %v", err, test.wantError)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
				t.Errorf("%+v != %+v", got, test.want)
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
			input: "uid",
			want: &property{
				key:  "uid",
				vals: map[string][]string{},
			},
		},
		{
			input:       "",
			want:        nil,
			shouldError: true,
			wantError:   &invalidKeyErr{"property key cannot be empty"},
		},
	}

	for _, test := range tests {
		got, err := sec.get(test.input)

		if test.shouldError {
			if !reflect.DeepEqual(err, test.wantError) {
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

func TestPropertyAdd(t *testing.T) {
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
					"en": []string{"Hello"},
					"fr": []string{"Bonjour"},
				},
			},
		},
	}

	for _, test := range tests {
		got := newProperty(test.key)

		for k, v := range test.values {
			for _, vv := range v {
				got.add(k, vv)
			}
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}

func TestPropertyGet(t *testing.T) {
	prop := property{
		key: "shell",
		vals: map[string][]string{
			"":      []string{"/bin/bash"},
			"win32": []string{"PowerShell.exe"},
		},
	}

	tests := []struct {
		desc  string
		input string
		want  []string
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
			desc:  "create",
			input: "haiku",
			want:  []string{},
		},
	}

	for _, test := range tests {
		got := prop.get(test.input)

		if !cmp.Equal(got, test.want) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}
