package ini

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type point struct {
	x, y int
}

func (p point) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("(%v,%v)", p.x, p.y)), nil
}

func TestEncodeProperty(t *testing.T) {
	tests := []struct {
		desc  string
		input struct {
			key string
			val interface{}
		}
		want        *bytes.Buffer
		shouldError bool
		wantError   error
	}{
		{
			desc: "encode string",
			input: struct {
				key string
				val interface{}
			}{"k", "0"},
			want: bytes.NewBufferString("k=0\n"),
		},
		{
			desc: "encode int",
			input: struct {
				key string
				val interface{}
			}{"k", 1},
			want: bytes.NewBufferString("k=1\n"),
		},
		{
			desc: "encode float",
			input: struct {
				key string
				val interface{}
			}{"k", 1.1},
			want: bytes.NewBufferString("k=1.1\n"),
		},
		{
			desc: "encode uint",
			input: struct {
				key string
				val interface{}
			}{"k", uint(1)},
			want: bytes.NewBufferString("k=1\n"),
		},
		{
			desc: "encode bool",
			input: struct {
				key string
				val interface{}
			}{"k", true},
			want: bytes.NewBufferString("k=true\n"),
		},
		{
			desc: "encode slice",
			input: struct {
				key string
				val interface{}
			}{"k", []string{"a", "b"}},
			want: bytes.NewBufferString("k=a\nk=b\n"),
		},
		{
			desc: "encode map",
			input: struct {
				key string
				val interface{}
			}{"k", map[string]string{"a": "a", "b": "b"}},
			want: bytes.NewBufferString("k[a]=a\nk[b]=b\n"),
		},
		{
			desc: "encode error map struct",
			input: struct {
				key string
				val interface{}
			}{"k", map[string]struct{}{"a": {}}},
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf(struct{}{})},
		},
		{
			desc: "encode struct",
			input: struct {
				key string
				val interface{}
			}{"", struct{}{}},
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf(struct{}{})},
		},
		{
			desc: "encode slice of ",
			input: struct {
				key string
				val interface{}
			}{"", []struct{}{{}}},
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf([]struct{}{})},
		},
		{
			desc: "encode MarshalText",
			input: struct {
				key string
				val interface{}
			}{"k", point{1, 3}},
			want: bytes.NewBufferString("k=(1,3)\n"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := new(bytes.Buffer)
			err := encodeProperty(got, test.input.key, reflect.ValueOf(test.input.val))

			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(MarshalTypeError{})) {
					t.Fatalf("encodeProperty(%v, %#v) returned %v, want %v", test.input.key, test.input.val, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("encodeProperty(%v, %#v) returned %v, want %v", test.input.key, test.input.val, err, test.wantError)
				}

				if !cmp.Equal(got, test.want, cmp.AllowUnexported(bytes.Buffer{})) {
					t.Errorf("encodeProperty(%v, %#v) = %#v, want %#v", test.input.key, test.input.val, string(got.Bytes()), string(test.want.Bytes()))
				}
			}
		})
	}
}

func TestEncodeSection(t *testing.T) {
	tests := []struct {
		desc  string
		input struct {
			key string
			val interface{}
		}
		want        *bytes.Buffer
		shouldError bool
		wantError   error
	}{
		{
			desc: "omit tag explicitly",
			input: struct {
				key string
				val interface{}
			}{"s", struct {
				P string `ini:"-"`
			}{"v"}},
			want: bytes.NewBufferString("[s]\n\n"),
		},
		{
			desc: "omitempty omits zero value",
			input: struct {
				key string
				val interface{}
			}{"s", struct {
				Z int `ini:",omitempty"`
				N int `ini:"N"`
			}{0, 1}},
			want: bytes.NewBufferString("[s]\nN=1\n\n"),
		},
		{
			desc: "encode error non-struct",
			input: struct {
				key string
				val interface{}
			}{"s", "string"},
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf("")},
		},
		{
			desc: "encode error property",
			input: struct {
				key string
				val interface{}
			}{"s", struct{ P struct{} }{}},
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf(struct{}{})},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := new(bytes.Buffer)
			err := encodeSection(got, test.input.key, reflect.ValueOf(test.input.val))

			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(MarshalTypeError{})) {
					t.Fatalf("encodeSection(%v, %#v) returned %v, want %v", test.input.key, test.input.val, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("encodeSection(%v, %#v) returned %v, want %v", test.input.key, test.input.val, err, test.wantError)
				}

				if !cmp.Equal(got, test.want, cmp.AllowUnexported(bytes.Buffer{})) {
					t.Errorf("encodeSection(%v, %#v) = %#v, want %#v", test.input.key, test.input.val, string(got.Bytes()), string(test.want.Bytes()))
				}
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		desc        string
		input       interface{}
		want        *bytes.Buffer
		shouldError bool
		wantError   error
	}{
		{
			desc:        "encode error not struct",
			input:       "",
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf("")},
		},
		{
			desc: "omitempty omits zero-value field",
			input: struct {
				Z int `ini:",omitempty"`
				N int
			}{0, 1},
			want: bytes.NewBufferString("N=1\n\n"),
		},
		{
			desc:        "encode error property",
			input:       struct{ P complex64 }{},
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf(complex64(10))},
		},
		{
			desc: "encode rv pointer",
			input: func() interface{} {
				return &struct{}{}
			}(),
			want: bytes.NewBufferString("\n"),
		},
		{
			desc: "omitempty omits zero-value struct field",
			input: struct {
				Z struct{} `ini:",omitempty"`
				N struct{}
			}{struct{}{}, struct{}{}},
			want: bytes.NewBufferString("\n[N]\n\n"),
		},
		{
			desc:        "encode error section property",
			input:       struct{ S struct{ P struct{} } }{},
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf(struct{}{})},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := new(bytes.Buffer)
			err := encode(got, reflect.ValueOf(test.input))

			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(MarshalTypeError{})) {
					t.Fatalf("encode(%#v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("encode(%#v) returned %v, want %v", test.input, err, test.wantError)
				}

				if !cmp.Equal(got, test.want, cmp.AllowUnexported(bytes.Buffer{})) {
					t.Errorf("encode(%#v) = %#v, want %#v", test.input, string(got.Bytes()), string(test.want.Bytes()))
				}
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		desc        string
		input       interface{}
		want        []byte
		shouldError bool
		wantError   error
	}{
		{
			desc:        "marshal error",
			input:       struct{ S struct{ P struct{} } }{},
			want:        nil,
			shouldError: true,
			wantError:   &MarshalTypeError{reflect.TypeOf(struct{}{})},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := Marshal(test.input)
			if test.shouldError {
				if !cmp.Equal(err, test.wantError, cmpopts.IgnoreUnexported(MarshalTypeError{})) {
					t.Fatalf("Marshal(%#v) returned %v, want %v", test.input, err, test.wantError)
				}
			} else {
				if err != nil {
					t.Fatalf("Marshal(%#v) returned %v, want %v", test.input, err, test.wantError)
				}

				if !cmp.Equal(got, test.want) {
					t.Errorf("Marshal(%#v) = %v, want %#v", test.input, string(got), string(test.want))
				}
			}
		})
	}
}
