package ini

import (
	"reflect"
	"regexp"
	"strings"
)

type tag struct {
	name      string
	omitempty bool
}

func newTag(sf reflect.StructField) tag {
	var t tag
	st := strings.SplitN(sf.Tag.Get("ini"), ",", 2)
	switch len(st) {
	case 1:
		t.name = st[0]
		if t.name == "" {
			t.name = sf.Name
		}
	default:
		t.name = st[0]
		t.omitempty = strings.Contains(st[1], "omitempty")
	}
	return t
}

func (t tag) pattern() (*regexp.Regexp, error) {
	if strings.HasPrefix(t.name, "[") && strings.HasSuffix(t.name, "]") {
		return regexp.Compile(t.name[1 : len(t.name)-1])
	}
	return nil, nil
}
