package ini

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input       interface{}
		want        string
		shouldError bool
		wantError   error
	}{
		{
			input: "/bin/bash",
			want:  "/bin/bash",
		},
		{
			input:       42,
			want:        "",
			shouldError: true,
			wantError: &UnmarshalTypeError{
				Value: reflect.ValueOf(42).String(),
				Type:  reflect.PtrTo(reflect.TypeOf("")),
			},
		},
	}

	for _, test := range tests {
		var got string
		rv := reflect.ValueOf(&got)

		err := decodeString(test.input, rv)

		if test.shouldError {
			if !reflect.DeepEqual(err, test.wantError) {
				t.Errorf("%v != %v", err, test.wantError)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("%v != %v", got, test.want)
			}
		}
	}
}

func TestDecodeInt(t *testing.T) {
	tests := []struct {
		input       string
		want        int64
		shouldError bool
		wantError   error
	}{
		{
			input: "42",
			want:  int64(42),
		},
		{
			input:       "forty-two",
			want:        int64(42),
			shouldError: true,
			wantError: &UnmarshalTypeError{
				Value: reflect.ValueOf("forty-two").String(),
				Type:  reflect.PtrTo(reflect.TypeOf(int64(42))),
			},
		},
	}

	for _, test := range tests {
		var got int64
		rv := reflect.ValueOf(&got)

		err := decodeInt(test.input, rv)

		if test.shouldError {
			if !reflect.DeepEqual(err, test.wantError) {
				t.Errorf("%v != %v", err, test.wantError)
			}
		} else {
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("%v != %v", got, test.want)
			}
		}
	}
}

func TestDecodeStruct(t *testing.T) {
	type user struct {
		Shell  string   `ini:"shell"`
		UID    int      `ini:"uid"`
		Groups []string `ini:"group"`
	}
	tests := []struct {
		input       section
		want        user
		shouldError bool
		wantError   error
	}{
		{
			input: section{
				name: "user",
				props: map[string]property{
					"shell": property{
						key: "shell",
						val: []string{"/bin/bash"},
					},
					"uid": property{
						key: "uid",
						val: []string{"1000"},
					},
					"group": property{
						key: "group",
						val: []string{"wheel", "video"},
					},
				},
			},
			want: user{
				Shell:  "/bin/bash",
				UID:    1000,
				Groups: []string{"wheel", "video"},
			},
		},
	}

	for _, test := range tests {
		var got user
		rv := reflect.ValueOf(&got)

		err := decodeStruct(test.input, rv)

		if test.shouldError {
			if !reflect.DeepEqual(err, test.wantError) {
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

func TestDecodeSlice(t *testing.T) {
	var tests []struct {
		input       property
		want        interface{}
		shouldError bool
		wantError   error
	}

	/*** []string tests ***/
	tests = []struct {
		input       property
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: property{
				key: "",
				val: []string{"/bin/bash", "/bin/zsh"},
			},
			want: []string{"/bin/bash", "/bin/zsh"},
		},
	}

	for _, test := range tests {
		var got []string

		err := decodeSlice(test.input, reflect.ValueOf(&got))

		if test.shouldError {
			if !reflect.DeepEqual(err, test.wantError) {
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

	/*** []int tests ***/
	tests = []struct {
		input       property
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: property{
				key: "",
				val: []string{"1000", "1001"},
			},
			want: []int{1000, 1001},
		},
	}

	for _, test := range tests {
		var got []int

		err := decodeSlice(test.input, reflect.ValueOf(&got))

		if test.shouldError {
			if !reflect.DeepEqual(err, test.wantError) {
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

func TestDecode(t *testing.T) {
	type user struct {
		Shell  string   `ini:"shell"`
		UID    int      `ini:"uid"`
		Groups []string `ini:"group"`
	}
	type config struct {
		User    user     `ini:"user"`
		Sources []string `ini:"source"`
	}

	tests := []struct {
		input ast
		want  config
	}{
		{
			input: ast{
				"": []section{
					section{
						name: "",
						props: map[string]property{
							"source": property{
								key: "source",
								val: []string{"passwd", "ldap"},
							},
						},
					},
				},
				"user": []section{
					section{
						name: "user",
						props: map[string]property{
							"shell": property{
								key: "shell",
								val: []string{"/bin/bash"},
							},
							"uid": property{
								key: "uid",
								val: []string{"42"},
							},
							"group": property{
								key: "group",
								val: []string{"wheel", "video"},
							},
						},
					},
				},
			},
			want: config{
				User: user{
					Shell:  "/bin/bash",
					UID:    42,
					Groups: []string{"wheel", "video"},
				},
				Sources: []string{"passwd", "ldap"},
			},
		},
	}

	for _, test := range tests {
		var got config
		rv := reflect.ValueOf(&got)

		err := decode(test.input, rv)

		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(got, test.want) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}

func TestUnmarshal(t *testing.T) {
	type user struct {
		Shell  string   `ini:"shell"`
		UID    int      `ini:"uid"`
		Groups []string `ini:"group"`
	}
	type config struct {
		User    user     `ini:"user"`
		Sources []string `ini:"source"`
	}

	tests := []struct {
		input string
		want  config
	}{
		{
			input: `source=passwd
[user]
shell=/bin/bash
uid=1000
group=wheel
group=video`,
			want: config{
				Sources: []string{"passwd"},
				User: user{
					Shell:  "/bin/bash",
					UID:    1000,
					Groups: []string{"wheel", "video"},
				},
			},
		},
	}

	for _, test := range tests {
		var got config
		err := Unmarshal([]byte(test.input), &got)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(got, test.want) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}
