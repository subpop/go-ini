package ini

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeString(t *testing.T) {
	tests := []struct {
		description string
		input       string
		want        string
		shouldError bool
		wantError   error
	}{
		{
			description: "valid",
			input:       "/bin/bash",
			want:        "/bin/bash",
		},
	}

	for _, test := range tests {
		var got string
		rv := reflect.ValueOf(&got)

		err := decodeString(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{})) {
				t.Fatalf("decodeString(%v) returned %v, want %v", test.input, err, test.wantError)
			}
		} else {
			if err != nil {
				t.Fatalf("decodeString(%v) returned %v, want %v", test.input, err, test.wantError)
			}
			if !cmp.Equal(got, test.want) {
				t.Errorf("decodeString(%v) = %v, want %v", test.input, got, test.want)
			}
		}
	}
}

func TestDecodeInt(t *testing.T) {
	tests := []struct {
		description string
		input       string
		want        int64
		shouldError bool
		wantError   error
	}{
		{
			description: "valid",
			input:       "42",
			want:        int64(42),
		},
		{
			description: "invalid parse syntax",
			input:       "forty-two",
			shouldError: true,
			wantError:   &DecodeError{errors.New("invalid syntax")},
		},
	}

	for _, test := range tests {
		var got int64
		rv := reflect.ValueOf(&got)

		err := decodeInt(test.input, rv)

		if test.shouldError {
			if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{})) {
				t.Fatalf("decodeInt(%v) returned %v, want %v", test.input, err, test.wantError)
			}
		} else {
			if err != nil {
				t.Fatalf("decodeInt(%v) returned %v, want %v", test.input, err, test.wantError)
			}
			if !cmp.Equal(got, test.want) {
				t.Errorf("decodeInt(%v) = %v, want %v", test.input, got, test.want)
			}
		}
	}
}

func TestDecodeUint(t *testing.T) {
	tests := []struct {
		description string
		input       string
		want        uint64
		shouldError bool
		wantError   error
	}{
		{
			description: "valid",
			input:       "42",
			want:        uint64(42),
		},
		{
			description: "invalid parse syntax",
			input:       "forty-two",
			shouldError: true,
			wantError:   &DecodeError{errors.New("invalid syntax")},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var got uint64
			rv := reflect.ValueOf(&got)

			err := decodeUint(test.input, rv)

			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{})) {
					t.Fatalf("decodeUint(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("decodeUint(%v) returned %v, want %v", test.input, err, test.wantError)
				}
				if !cmp.Equal(got, test.want) {
					t.Errorf("decodeUint(%v) = %v, want %v", test.input, got, test.want)
				}
			}
		})
	}
}

func TestDecodeBool(t *testing.T) {
	tests := []struct {
		description string
		input       string
		want        bool
		shouldError bool
		wantError   error
	}{
		{
			description: "valid (true)",
			input:       "true",
			want:        true,
		},
		{
			description: "valid (0)",
			input:       "0",
			want:        false,
		},
		{
			description: "valid (T)",
			input:       "T",
			want:        true,
		},
		{
			description: "invalid parse syntax",
			input:       "forty-two",
			shouldError: true,
			wantError:   &DecodeError{errors.New("invalid syntax")},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var got bool
			rv := reflect.ValueOf(&got)

			err := decodeBool(test.input, rv)

			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{})) {
					t.Fatalf("decodeBool(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("decodeBool(%v) returned %v, want %v", test.input, err, test.wantError)
				}
				if !cmp.Equal(got, test.want) {
					t.Errorf("decodeBool(%v) = %v, want %v", test.input, got, test.want)
				}
			}
		})
	}
}

func TestDecodeFloat(t *testing.T) {
	tests := []struct {
		description string
		input       string
		want        float64
		shouldError bool
		wantError   error
	}{
		{
			description: "valid",
			input:       "42.2",
			want:        float64(42.2),
		},
		{
			description: "invalid parse syntax",
			input:       "forty-two",
			want:        float64(0),
			shouldError: true,
			wantError:   &DecodeError{errors.New("invalid syntax")},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			var got float64
			rv := reflect.ValueOf(&got)

			err := decodeFloat(test.input, rv)

			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{})) {
					t.Fatalf("decodeFloat(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("decodeFloat(%v) returned %v, want %v", test.input, err, test.wantError)
				}
				if !cmp.Equal(got, test.want) {
					t.Errorf("decodeFloat(%v) = %v, want %v", test.input, got, test.want)
				}
			}
		})
	}
}

func TestDecodeStruct(t *testing.T) {
	tests := []struct {
		description string
		input       section
		want        interface{}
		shouldError bool
		wantError   error
		init        func() interface{}
	}{
		{
			description: "decodeString",
			input:       section{"section", map[string]property{"property": {"property", map[string][]string{"": {"value"}}}}},
			want: &struct {
				Property string `ini:"property"`
			}{"value"},
			init: func() interface{} {
				return &struct {
					Property string `ini:"property"`
				}{}
			},
		},
		{
			description: "decodeInt",
			input:       section{"section", map[string]property{"property": {"property", map[string][]string{"": {"0"}}}}},
			want: &struct {
				Property int `ini:"property"`
			}{0},
			init: func() interface{} {
				return &struct {
					Property int `ini:"property"`
				}{}
			},
		},
		{
			description: "decodeUint",
			input:       section{"section", map[string]property{"property": {"property", map[string][]string{"": {"0"}}}}},
			want: &struct {
				Property uint `ini:"property"`
			}{0},
			init: func() interface{} {
				return &struct {
					Property uint `ini:"property"`
				}{}
			},
		},
		{
			description: "decodeFloat",
			input:       section{"section", map[string]property{"property": {"property", map[string][]string{"": {"0.0"}}}}},
			want: &struct {
				Property float64 `ini:"property"`
			}{0.0},
			init: func() interface{} {
				return &struct {
					Property float64 `ini:"property"`
				}{}
			},
		},
		{
			description: "decodeBool",
			input:       section{"section", map[string]property{"property": {"property", map[string][]string{"": {"1"}}}}},
			want: &struct {
				Property bool `ini:"property"`
			}{true},
			init: func() interface{} {
				return &struct {
					Property bool `ini:"property"`
				}{}
			},
		},
		{
			description: "skip property",
			input:       section{"section", map[string]property{"property": {"property", map[string][]string{"": {"0"}}}}},
			want: &struct {
				Property int `ini:"-"`
			}{0},
			init: func() interface{} {
				return &struct {
					Property int `ini:"-"`
				}{}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := test.init()

			err := decodeStruct(test.input, reflect.ValueOf(got))
			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{}, UnmarshalTypeError{})) {
					t.Fatalf("decodeStruct(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("decodeStruct(%v) returned %v, want %v", test.input, err, test.wantError)
				}

				if !cmp.Equal(got, test.want) {
					t.Errorf("decodeStruct(%v) = %v, want %v", test.input, got, test.want)
				}
			}
		})
	}
}

func TestDecodeSliceStruct(t *testing.T) {
	tests := []struct {
		description string
		input       []section
		want        interface{}
		shouldError bool
		wantError   error
		init        func() interface{}
	}{
		{
			input: []section{
				{
					name: "section",
					props: map[string]property{
						"property": {
							key: "property",
							vals: map[string][]string{
								"": {"value0"},
							},
						},
					},
				},
				{
					name: "section",
					props: map[string]property{
						"property": {
							key: "property",
							vals: map[string][]string{
								"": {"value1"},
							},
						},
					},
				},
			},
			want: &[]struct {
				Property string `ini:"property"`
			}{
				{
					Property: "value0",
				},
				{
					Property: "value1",
				},
			},
			init: func() interface{} {
				return &[]struct {
					Property string `ini:"property"`
				}{}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := test.init()

			err := decodeSliceStruct(test.input, reflect.ValueOf(got))
			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{}, UnmarshalTypeError{})) {
					t.Fatalf("decodeSliceStruct(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("decodeSliceStruct(%v) returned %v, want %v", test.input, err, test.wantError)
				}

				if !cmp.Equal(got, test.want) {
					t.Errorf("decodeSliceStruct(%v) = %v, want %v", test.input, got, test.want)
				}
			}
		})
	}
}

func TestDecodeSlice(t *testing.T) {
	tests := []struct {
		description string
		input       []string
		want        interface{}
		shouldError bool
		wantError   error
		init        func() interface{}
	}{
		{
			description: "string slice",
			input:       []string{"value0", "value1"},
			want:        &[]string{"value0", "value1"},
			init: func() interface{} {
				return &[]string{}
			},
		},
		{
			description: "int slice",
			input:       []string{"0", "1"},
			want:        &[]int{0, 1},
			init: func() interface{} {
				return &[]int{}
			},
		},
		{
			description: "uint slice",
			input:       []string{"0", "1"},
			want:        &[]uint{0, 1},
			init: func() interface{} {
				return &[]uint{}
			},
		},
		{
			description: "float64 slice",
			input:       []string{"0.0", "1.0"},
			want:        &[]float64{0.0, 1.0},
			init: func() interface{} {
				return &[]float64{}
			},
		},
		{
			description: "bool slice",
			input:       []string{"0", "1"},
			want:        &[]bool{false, true},
			init: func() interface{} {
				return &[]bool{}
			},
		},
		{
			description: "struct slice",
			input:       []string{"0", "1"},
			want:        &[]struct{}{{}, {}},
			shouldError: true,
			wantError: &UnmarshalTypeError{
				val: reflect.ValueOf([]string{"0", "1"}).String(),
				typ: reflect.PtrTo(reflect.TypeOf([]struct{}{})),
			},
			init: func() interface{} {
				return &[]struct{}{}
			},
		},
		{
			description: "invalid parse syntax",
			input:       []string{"one", "two"},
			want:        &[]bool{true, false},
			shouldError: true,
			wantError:   &DecodeError{errors.New("invalid syntax")},
			init: func() interface{} {
				return &[]bool{}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := test.init()

			err := decodeSlice(test.input, reflect.ValueOf(got))
			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{}, UnmarshalTypeError{})) {
					t.Fatalf("decodeSlice(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("decodeSlice(%v) returned %v, want %v", test.input, err, test.wantError)
				}

				if !cmp.Equal(got, test.want) {
					t.Errorf("decodeSlice(%v) = %v, want %v", test.input, got, test.want)
				}
			}
		})
	}
}

func TestDecodeMap(t *testing.T) {
	tests := []struct {
		description string
		input       property
		want        interface{}
		shouldError bool
		wantError   error
		init        func() interface{}
	}{
		{
			description: "map[string]string",
			input:       property{key: "p", vals: map[string][]string{"k1": {"v1"}, "k2": {"v2"}}},
			want:        &map[string]string{"k1": "v1", "k2": "v2"},
			init: func() interface{} {
				return &map[string]string{}
			},
		},
		{
			description: "map[string]int",
			input:       property{key: "p", vals: map[string][]string{"k1": {"0"}, "k2": {"1"}}},
			want:        &map[string]int{"k1": 0, "k2": 1},
			init: func() interface{} {
				return &map[string]int{}
			},
		},
		{
			description: "map[string]uint",
			input:       property{key: "p", vals: map[string][]string{"k1": {"0"}, "k2": {"1"}}},
			want:        &map[string]uint{"k1": 0, "k2": 1},
			init: func() interface{} {
				return &map[string]uint{}
			},
		},
		{
			description: "map[string]float64",
			input:       property{key: "p", vals: map[string][]string{"k1": {"0.0"}, "k2": {"1.0"}}},
			want:        &map[string]float64{"k1": 0.0, "k2": 1.0},
			init: func() interface{} {
				return &map[string]float64{}
			},
		},
		{
			description: "map[string]bool",
			input:       property{key: "p", vals: map[string][]string{"k1": {"0"}, "k2": {"1"}}},
			want:        &map[string]bool{"k1": false, "k2": true},
			init: func() interface{} {
				return &map[string]bool{}
			},
		},
		{
			description: "map[string][]string",
			input:       property{key: "p", vals: map[string][]string{"k1": {"v0", "v1"}}},
			want:        &map[string][]string{"k1": {"v0", "v1"}},
			init: func() interface{} {
				return &map[string][]string{}
			},
		},
		{
			description: "map[string]struct{}",
			input:       property{key: "p", vals: map[string][]string{"k1": {"0"}, "k2": {"1"}}},
			want:        &map[string]struct{}{},
			shouldError: true,
			wantError:   &UnmarshalTypeError{val: reflect.ValueOf(property{}).String(), typ: reflect.TypeOf(&map[string]struct{}{})},
			init: func() interface{} {
				return &map[string]struct{}{}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := test.init()

			err := decodeMap(test.input, reflect.ValueOf(got))
			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{}, UnmarshalTypeError{})) {
					t.Fatalf("decodeMap(%v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("decodeMap(%v) returned %v, want %v", test.input, err, test.wantError)
				}

				if !cmp.Equal(got, test.want) {
					t.Errorf("decodeMap(%v) = %v, want %v", test.input, got, test.want)
				}
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		description string
		input       parseTree
		want        interface{}
		shouldError bool
		wantError   error
		init        func() interface{}
	}{
		{
			description: "decode top-level struct",
			input: parseTree{
				global: section{
					name: "",
					props: map[string]property{
						"property": {
							key: "property",
							vals: map[string][]string{
								"": {"value"},
							},
						},
						"map": {
							key: "map",
							vals: map[string][]string{
								"k1": {"v1"},
								"k2": {"v2"},
							},
						},
					},
				},
				sections: map[string][]section{
					"section1": {
						{
							name: "section1",
							props: map[string]property{
								"key1": {
									key: "key1",
									vals: map[string][]string{
										"": {"value1"},
									},
								},
							},
						},
						{
							name: "section1",
							props: map[string]property{
								"key1": {
									key: "key1",
									vals: map[string][]string{
										"": {"value2"},
									},
								},
							},
						},
					},
				},
			},
			want: &struct {
				Property string            `ini:"property"`
				Map      map[string]string `ini:"map"`
				Section  []struct {
					Key string `ini:"key1"`
				} `ini:"section1"`
			}{
				Property: "value",
				Map: map[string]string{
					"k1": "v1",
					"k2": "v2",
				},
				Section: []struct {
					Key string `ini:"key1"`
				}{
					{Key: "value1"},
					{Key: "value2"},
				},
			},
			init: func() interface{} {
				return &struct {
					Property string            `ini:"property"`
					Map      map[string]string `ini:"map"`
					Section  []struct {
						Key string `ini:"key1"`
					} `ini:"section1"`
				}{}
			},
		},
		{
			description: "decode top-level map",
			input:       parseTree{},
			want:        map[string]interface{}{},
			init: func() interface{} {
				return map[string]interface{}{}
			},
			shouldError: true,
			wantError:   &DecodeError{err: fmt.Errorf("cannot unmarshal into value of type map[string]interface {}")},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := test.init()

			err := decode(test.input, reflect.ValueOf(got))
			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(DecodeError{}, UnmarshalTypeError{})) {
					t.Fatalf("decode(%+v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("decode(%+v) returned %v, want %v", test.input, err, test.wantError)
				}

				if !cmp.Equal(got, test.want) {
					t.Errorf("decode(%+v) = %+v, want %+v\ndiff -want +got\n%v", test.input, got, test.want, cmp.Diff(test.want, got))
				}
			}
		})
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
		Version string   `ini:"version"`
	}

	tests := []struct {
		description string
		input       string
		want        config
		shouldError bool
		wantError   error
	}{
		{
			input: `source=passwd
kver=1

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
		t.Run(test.description, func(t *testing.T) {
			var got config
			err := Unmarshal([]byte(test.input), &got)
			if err != nil {
				t.Fatalf("Unmarshal(%v) returned %v, want %v", test.input, err, test.wantError)
			}

			if !cmp.Equal(got, test.want) {
				t.Errorf("Unmarshal(%v) = %v, want %v\ndiff -want +got\n%v", test.input, got, want, cmp.Diff(test.want, got))
			}
		})
	}
}

func TestUnmarshalWildcard(t *testing.T) {
	type User struct {
		ININame string
		Shell   string `ini:"shell"`
	}
	tests := []struct {
		description string
		input       string
		want        struct {
			Users []User `ini:"*"`
		}
		shouldError bool
		wantError   error
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
		t.Run(test.description, func(t *testing.T) {
			var got struct {
				Users []User `ini:"*"`
			}
			err := Unmarshal([]byte(test.input), &got)
			if err != nil {
				t.Fatalf("Unmarshal(%+v) return %v, want %v", test.input, err, test.wantError)
			}

			if !cmp.Equal(got, test.want, cmp.Options{cmpopts.SortSlices(func(a, b User) bool {
				return a.ININame > b.ININame
			})}) {
				t.Errorf("Unmarshal(%+v) = %v, want %v\ndiff -want +got\n%v", test.input, got, test.want, cmp.Diff(test.want, got))
			}
		})
	}
}
