package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAST(t *testing.T) {
	tests := []struct {
		sections []section
		want     ast
	}{
		{
			sections: []section{},
			want: ast{
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
							val: []string{"/bin/bash"},
						},
					},
				},
			},
			want: ast{
				"": []section{
					{name: "", props: map[string]property{}},
				},
				"user": []section{
					{
						name: "user",
						props: map[string]property{
							"shell": {
								key: "shell",
								val: []string{"/bin/bash"},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		got := newAST()

		for _, s := range test.sections {
			got.addSection(s)
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}

func TestSection(t *testing.T) {
	tests := []struct {
		name  string
		props []property
		want  section
	}{
		{
			name: "user",
			props: []property{
				{key: "name", val: []string{"root"}},
				{key: "shell", val: []string{"/bin/bash"}},
				{key: "uid", val: []string{"1000", "1001"}},
			},
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
					"uid": property{
						key: "uid",
						val: []string{"1000", "1001"},
					},
				},
			},
		},
	}

	for _, test := range tests {
		got := newSection(test.name)

		for _, p := range test.props {
			got.addProperty(p)
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{}, section{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}

func TestProperty(t *testing.T) {
	tests := []struct {
		key  string
		val  []string
		want property
	}{
		{
			key: "shell",
			val: []string{"/bin/bash", "/bin/zsh"},
			want: property{
				key: "shell",
				val: []string{"/bin/bash", "/bin/zsh"},
			},
		},
	}

	for _, test := range tests {
		got := newProperty(test.key)

		for _, v := range test.val {
			got.appendVal(v)
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmp.AllowUnexported(property{})}) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}
