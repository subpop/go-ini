package ini

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type advanced struct {
	Name    string `ini:"name"`
	Age     int    `ini:"age"`
	Address struct {
		Street string `ini:"street"`
		City   string `ini:"city"`
		State  string `ini:"state"`
		ZIP    string `ini:"zip"`
	} `ini:"address"`
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		input string
		want  advanced
	}{
		{
			input: `name=Rupert
			age=23
			
			[address]
			street=123 Main St.
			city=Boston
			state=Massachusetts
			zip=02108`,
			want: advanced{
				Name: "Rupert",
				Age:  23,
				Address: struct {
					Street string `ini:"street"`
					City   string `ini:"city"`
					State  string `ini:"state"`
					ZIP    string `ini:"zip"`
				}{
					Street: "123 Main St.",
					City:   "Boston",
					State:  "Massachusetts",
					ZIP:    "02108",
				},
			},
		},
	}

	for _, test := range tests {
		var got advanced
		if err := Unmarshal([]byte(test.input), &got); err != nil {
			t.Error(err)
		}

		if !cmp.Equal(got, test.want) {
			t.Errorf("%v != %v", got, test.want)
		}
	}
}
