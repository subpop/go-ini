package ini

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

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
				val: reflect.ValueOf(42).String(),
				typ: reflect.PtrTo(reflect.TypeOf("")),
			},
		},
	}

	for _, test := range tests {
		var got string
		rv := reflect.ValueOf(&got)

		err := decodeString(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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
				val: reflect.ValueOf("forty-two").String(),
				typ: reflect.PtrTo(reflect.TypeOf(int64(42))),
			},
		},
	}

	for _, test := range tests {
		var got int64
		rv := reflect.ValueOf(&got)

		err := decodeInt(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

func TestDecodeUint(t *testing.T) {
	tests := []struct {
		input       string
		want        uint64
		shouldError bool
		wantError   error
	}{
		{
			input: "42",
			want:  uint64(42),
		},
		{
			input:       "forty-two",
			want:        uint64(42),
			shouldError: true,
			wantError: &UnmarshalTypeError{
				val: reflect.ValueOf("forty-two").String(),
				typ: reflect.PtrTo(reflect.TypeOf(uint64(42))),
			},
		},
	}

	for _, test := range tests {
		var got uint64
		rv := reflect.ValueOf(&got)

		err := decodeUint(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

func TestDecodeBool(t *testing.T) {
	tests := []struct {
		input       string
		want        bool
		shouldError bool
		wantError   error
	}{
		{
			input: "true",
			want:  true,
		},
		{
			input: "0",
			want:  false,
		},
		{
			input: "T",
			want:  true,
		},
		{
			input:       "forty-two",
			want:        false,
			shouldError: true,
			wantError: &UnmarshalTypeError{
				val: reflect.ValueOf("forty-two").String(),
				typ: reflect.PtrTo(reflect.TypeOf(false)),
			},
		},
	}

	for _, test := range tests {
		var got bool
		rv := reflect.ValueOf(&got)

		err := decodeBool(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

func TestDecodeFloat(t *testing.T) {
	tests := []struct {
		input       string
		want        float64
		shouldError bool
		wantError   error
	}{
		{
			input: "42.2",
			want:  float64(42.2),
		},
		{
			input:       "forty-two",
			want:        float64(42.2),
			shouldError: true,
			wantError: &UnmarshalTypeError{
				val: reflect.ValueOf("forty-two").String(),
				typ: reflect.PtrTo(reflect.TypeOf(float64(42.2))),
			},
		},
	}

	for _, test := range tests {
		var got float64
		rv := reflect.ValueOf(&got)

		err := decodeFloat(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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
		Shell    string   `ini:"shell"`
		UID      int      `ini:"uid"`
		Groups   []string `ini:"group"`
		GID      uint     `ini:"gid"`
		Admin    bool     `ini:"admin"`
		AccessID float32  `ini:"access_id"`
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
					"shell": {
						key: "shell",
						vals: map[string][]string{
							"": {"/bin/bash"},
						},
					},
					"uid": {
						key: "uid",
						vals: map[string][]string{
							"": {"1000"},
						},
					},
					"group": {
						key: "group",
						vals: map[string][]string{
							"": {"wheel", "video"},
						},
					},
					"gid": {
						key: "gid",
						vals: map[string][]string{
							"": {"1000"},
						},
					},
					"admin": {
						key: "admin",
						vals: map[string][]string{
							"": {"true"},
						},
					},
					"access_id": {
						key: "access_id",
						vals: map[string][]string{
							"": {"1234.5678"},
						},
					},
				},
			},
			want: user{
				Shell:    "/bin/bash",
				UID:      1000,
				Groups:   []string{"wheel", "video"},
				GID:      1000,
				Admin:    true,
				AccessID: 1234.5678,
			},
		},
	}

	for _, test := range tests {
		var got user
		rv := reflect.ValueOf(&got)

		err := decodeStruct(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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
		input       interface{}
		want        interface{}
		shouldError bool
		wantError   error
	}

	/*** []string tests ***/
	tests = []struct {
		input       interface{}
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: []string{"/bin/bash", "/bin/zsh"},
			want:  []string{"/bin/bash", "/bin/zsh"},
		},
	}

	for _, test := range tests {
		var got []string

		err := decodeSlice(test.input, reflect.ValueOf(&got))

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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
		input       interface{}
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: []string{"1000", "1001"},
			want:  []int{1000, 1001},
		},
	}

	for _, test := range tests {
		var got []int

		err := decodeSlice(test.input, reflect.ValueOf(&got))

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

	/*** []struct tests ***/
	type user struct {
		Name  string `ini:"name"`
		Shell string `ini:"shell"`
	}
	tests = []struct {
		input       interface{}
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: []section{
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
								"": {"/bin/bash"},
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
								"": {"/bin/zsh"},
							},
						},
					},
				},
			},
			want: []user{
				{
					Name:  "root",
					Shell: "/bin/bash",
				},
				{
					Name:  "admin",
					Shell: "/bin/zsh",
				},
			},
		},
	}

	for _, test := range tests {
		var got []user

		err := decodeSlice(test.input.([]section), reflect.ValueOf(&got))

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

	/*** []bool tests ***/
	tests = []struct {
		input       interface{}
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: []string{"true", "false"},
			want:  []bool{true, false},
		},
	}

	for _, test := range tests {
		var got []bool

		err := decodeSlice(test.input, reflect.ValueOf(&got))

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

	/*** []float64 tests ***/
	tests = []struct {
		input       interface{}
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: []string{"123.456", "654.321"},
			want:  []float64{123.456, 654.321},
		},
	}

	for _, test := range tests {
		var got []float64

		err := decodeSlice(test.input, reflect.ValueOf(&got))

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

func TestDecodeMap(t *testing.T) {

	/* map[string]string tests */
	tests := []struct {
		input       property
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: property{
				key: "Greeting",
				vals: map[string][]string{
					"en": {"Hello"},
					"fr": {"Bonjour"},
				},
			},
			want: map[string]string{
				"en": "Hello",
				"fr": "Bonjour",
			},
		},
	}

	for _, test := range tests {
		var got map[string]string
		rv := reflect.ValueOf(&got)

		err := decodeMap(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

	/* map[string]int tests */
	tests = []struct {
		input       property
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: property{
				key: "Planets",
				vals: map[string][]string{
					"Mercury": {"1"},
					"Venus":   {"2"},
				},
			},
			want: map[string]int{
				"Mercury": 1,
				"Venus":   2,
			},
		},
	}

	for _, test := range tests {
		var got map[string]int
		rv := reflect.ValueOf(&got)

		err := decodeMap(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

	/* map[string]float64 tests */
	tests = []struct {
		input       property
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: property{
				key: "Currency",
				vals: map[string][]string{
					"USD": {"1.0"},
					"GBP": {"1.2"},
				},
			},
			want: map[string]float64{
				"USD": 1.0,
				"GBP": 1.2,
			},
		},
	}

	for _, test := range tests {
		var got map[string]float64
		rv := reflect.ValueOf(&got)

		err := decodeMap(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

	/* map[string]bool tests */
	tests = []struct {
		input       property
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: property{
				key: "Switch",
				vals: map[string][]string{
					"On":  {"true"},
					"Off": {"false"},
				},
			},
			want: map[string]bool{
				"On":  true,
				"Off": false,
			},
		},
	}

	for _, test := range tests {
		var got map[string]bool
		rv := reflect.ValueOf(&got)

		err := decodeMap(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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

	/* map[string]uint tests */
	tests = []struct {
		input       property
		want        interface{}
		shouldError bool
		wantError   error
	}{
		{
			input: property{
				key: "Planets",
				vals: map[string][]string{
					"Mercury": {"1"},
					"Venus":   {"2"},
				},
			},
			want: map[string]uint{
				"Mercury": 1,
				"Venus":   2,
			},
		},
	}

	for _, test := range tests {
		var got map[string]uint
		rv := reflect.ValueOf(&got)

		err := decodeMap(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(UnmarshalTypeError{})) {
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
		input parseTree
		want  config
	}{
		{
			input: parseTree{
				global: section{
					name: "",
					props: map[string]property{
						"source": {
							key: "source",
							vals: map[string][]string{
								"": {"passwd", "ldap"},
							},
						},
					},
				},
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
								"uid": {
									key: "uid",
									vals: map[string][]string{
										"": {"42"},
									},
								},
								"group": {
									key: "group",
									vals: map[string][]string{
										"": {"wheel", "video"},
									},
								},
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
		Name   string            `ini:"name"`
		Shell  map[string]string `ini:"shell"`
		UID    int               `ini:"uid"`
		Groups []string          `ini:"group"`
	}
	type config struct {
		Users   []user   `ini:"user"`
		Sources []string `ini:"source"`
	}

	tests := []struct {
		input string
		want  config
	}{
		{
			input: `source=passwd
[user]
name=root
shell[unix]=/bin/bash
shell[win32]=PowerShell.exe
uid=1000
group=wheel
group=video

[user]
name=admin
shell[unix]=/bin/bash
shell[win32]=PowerShell.exe
uid=1001
group=wheel
group=video`,
			want: config{
				Sources: []string{"passwd"},
				Users: []user{
					{
						Name: "root",
						Shell: map[string]string{
							"unix":  "/bin/bash",
							"win32": "PowerShell.exe",
						},
						UID:    1000,
						Groups: []string{"wheel", "video"},
					},
					{
						Name: "admin",
						Shell: map[string]string{
							"unix":  "/bin/bash",
							"win32": "PowerShell.exe",
						},
						UID:    1001,
						Groups: []string{"wheel", "video"},
					},
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
			t.Errorf("%+v != %+v", got, test.want)
		}
	}
}

func TestUnmarshalWildcard(t *testing.T) {
	type User struct {
		ININame string
		Shell   string `ini:"shell"`
	}
	tests := []struct {
		input string
		want  struct {
			Users []User `ini:"*"`
		}
	}{
		{
			input: "[root]\nshell=/bin/bash\n\n[admin]\nshell=/bin/zsh",
			want: struct {
				Users []User `ini:"*"`
			}{
				Users: []User{
					{
						ININame: "root",
						Shell:   "/bin/bash",
					},
					{
						ININame: "admin",
						Shell:   "/bin/zsh",
					},
				},
			},
		},
	}

	for _, test := range tests {
		var got struct {
			Users []User `ini:"*"`
		}
		err := Unmarshal([]byte(test.input), &got)
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(got, test.want, cmp.Options{cmpopts.SortSlices(func(a, b User) bool {
			return a.ININame > b.ININame
		})}) {
			t.Errorf("%+v != %+v", got, test.want)
		}
	}
}
