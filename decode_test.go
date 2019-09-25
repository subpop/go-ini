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
				Value: reflect.ValueOf("forty-two").String(),
				Type:  reflect.PtrTo(reflect.TypeOf(uint64(42))),
			},
		},
	}

	for _, test := range tests {
		var got uint64
		rv := reflect.ValueOf(&got)

		err := decodeUint(test.input, rv)

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
				Value: reflect.ValueOf("forty-two").String(),
				Type:  reflect.PtrTo(reflect.TypeOf(false)),
			},
		},
	}

	for _, test := range tests {
		var got bool
		rv := reflect.ValueOf(&got)

		err := decodeBool(test.input, rv)

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
				Value: reflect.ValueOf("forty-two").String(),
				Type:  reflect.PtrTo(reflect.TypeOf(float64(42.2))),
			},
		},
	}

	for _, test := range tests {
		var got float64
		rv := reflect.ValueOf(&got)

		err := decodeFloat(test.input, rv)

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
			input: property{
				key: "",
				val: []string{"/bin/bash", "/bin/zsh"},
			},
			want: []string{"/bin/bash", "/bin/zsh"},
		},
	}

	for _, test := range tests {
		var got []string

		err := decodeSlice(test.input.(property).val, reflect.ValueOf(&got))

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
		input       interface{}
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

		err := decodeSlice(test.input.(property).val, reflect.ValueOf(&got))

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
				{
					name: "user",
					props: map[string]property{
						"name": property{
							key: "name",
							val: []string{"admin"},
						},
						"shell": property{
							key: "shell",
							val: []string{"/bin/zsh"},
						},
					},
				},
			},
			want: []user{
				user{
					Name:  "root",
					Shell: "/bin/bash",
				},
				user{
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
		Name   string   `ini:"name"`
		Shell  string   `ini:"shell"`
		UID    int      `ini:"uid"`
		Groups []string `ini:"group"`
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
shell=/bin/bash
uid=1000
group=wheel
group=video

[user]
name=admin
shell=/bin/bash
uid=1001
group=wheel
group=video`,
			want: config{
				Sources: []string{"passwd"},
				Users: []user{
					user{
						Name:   "root",
						Shell:  "/bin/bash",
						UID:    1000,
						Groups: []string{"wheel", "video"},
					},
					user{
						Name:   "admin",
						Shell:  "/bin/bash",
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

func TestUnmarshalPattern(t *testing.T) {
	type image struct {
		SectionName    string
		Name           string `ini:"name"`
		OSInfo         string `ini:"osinfo"`
		Arch           string `ini:"arch"`
		File           string `ini:"file"`
		Revision       int    `ini:"revision,omitempty"`
		Checksum       string `ini:"checksum"`
		Format         string `ini:"format"`
		Size           int64  `ini:"size"`
		CompressedSize int64  `ini:"compressed_size"`
		Expand         string `ini:"expand"`
		Notes          string `ini:"notes"`
	}
	type index struct {
		Images []image `ini:"[centos-.*]"`
	}
	tests := []struct {
		input string
		want  index
	}{
		{
			input: `[centos-6]
name=CentOS 6.6
osinfo=centos6.6
arch=x86_64
file=centos-6.xz
revision=6
checksum=fc403ea3555a5608a25ad30ce2514b67288311a7197ddf9fb664475820f26db2bd95a86be9cd6e3f772187b384a02e0965430456dd518d343a80457057bc5441
format=raw
size=6442450944
compressed_size=199265736
expand=/dev/sda3
notes=CentOS 6.6
	
	This CentOS image contains only unmodified @Core group packages.
	
	It is thus very minimal.  The kickstart and install script can be
	found in the libguestfs source tree:
	
	builder/website/centos.sh
	
	Note that ‘virt-builder centos-6’ will always install the latest
	6.x release.

[centos-7.0]
name=CentOS 7.0
osinfo=centos7.0
arch=x86_64
file=centos-7.0.xz
checksum=cf9ae295f633fbd04e575eeca16f372e933c70c3107c44eb06864760d04354aa94b4f356bfc9a598c138e687304a52e96777e4467e7db1ec0cb5b2d2ec61affc
format=raw
size=6442450944
compressed_size=213203844
expand=/dev/sda3
notes=CentOS 7.0
	
	This CentOS image contains only unmodified @Core group packages.
	
	It is thus very minimal.  The kickstart and install script can be
	found in the libguestfs source tree:
	
	builder/website/centos.sh

[debian-10]
name=Debian 10 (buster)
osinfo=debian10
arch=x86_64
file=debian-10.xz
checksum[sha512]=264d340e843d349f8caee14add56da4de95b22224ec48c6b3d9245afc764e4d460edabaf16fe6e4026008383128dc878a6d85eaf5dc55d66cef55cca88929c05
format=raw
size=6442450944
compressed_size=218919120
expand=/dev/sda1
notes=Debian 10 (buster)
	
	This is a minimal Debian install.
	
	This image is so very minimal that it only includes an ssh server
	This image does not contain SSH host keys.  To regenerate them use:
	
		--firstboot-command "dpkg-reconfigure openssh-server"
	
	This template was generated by a script in the libguestfs source tree:
		builder/templates/make-template.ml
	Associated files used to prepare this template can be found in the
	same directory.`,
			want: index{
				Images: []image{
					{
						SectionName:    "centos-6",
						Name:           "CentOS 6.6",
						OSInfo:         "centos6.6",
						Arch:           "x86_64",
						File:           "centos-6.xz",
						Revision:       6,
						Checksum:       "fc403ea3555a5608a25ad30ce2514b67288311a7197ddf9fb664475820f26db2bd95a86be9cd6e3f772187b384a02e0965430456dd518d343a80457057bc5441",
						Format:         "raw",
						Size:           6442450944,
						CompressedSize: 199265736,
						Expand:         "/dev/sda3",
						Notes:          "CentOS 6.6\n\t\n\tThis CentOS image contains only unmodified @Core group packages.\n\t\n\tIt is thus very minimal.  The kickstart and install script can be\n\tfound in the libguestfs source tree:\n\t\n\tbuilder/website/centos.sh\n\t\n\tNote that ‘virt-builder centos-6’ will always install the latest\n\t6.x release.",
					},
					{
						SectionName:    "centos-7.0",
						Name:           "CentOS 7.0",
						OSInfo:         "centos7.0",
						Arch:           "x86_64",
						File:           "centos-7.0.xz",
						Checksum:       "cf9ae295f633fbd04e575eeca16f372e933c70c3107c44eb06864760d04354aa94b4f356bfc9a598c138e687304a52e96777e4467e7db1ec0cb5b2d2ec61affc",
						Format:         "raw",
						Size:           6442450944,
						CompressedSize: 213203844,
						Expand:         "/dev/sda3",
						Notes:          "CentOS 7.0\n\t\n\tThis CentOS image contains only unmodified @Core group packages.\n\t\n\tIt is thus very minimal.  The kickstart and install script can be\n\tfound in the libguestfs source tree:\n\t\n\tbuilder/website/centos.sh",
					},
				},
			},
		},
	}

	for _, test := range tests {
		var got index
		if err := UnmarshalWithOptions([]byte(test.input), &got, Options{AllowMultilineValues: true}); err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(got, test.want) {
			t.Errorf("%+v != %+v", got, test.want)
		}
	}
}
