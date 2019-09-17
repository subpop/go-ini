package ini

import (
	"fmt"
	"reflect"
	"strconv"
)

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "ini: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "ini: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "ini: Unmarshal(nil " + e.Type.String() + ")"
}

func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	p := newParser(data)
	if err := p.parse(); err != nil {
		return fmt.Errorf("ini: decode: %v", err)
	}

	return decode(p, rv.Elem(), "")
}

func decode(p *parser, v reflect.Value, st string) error {
	if v.Kind() != reflect.Struct {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	for i := 0; i < v.NumField(); i++ {
		sf := v.Type().Field(i)
		tag, _ := key(sf)
		if tag == "-" {
			continue
		}
		switch sf.Type.Kind() {
		case reflect.String:
			sv := v.Field(i)
			sv.SetString(p.ast[st].props[tag].val[0])
		case reflect.Int:
			sv := v.Field(i)
			n, err := strconv.ParseInt(p.ast[st].props[tag].val[0], 10, 64)
			if err != nil {
				return err
			}
			sv.SetInt(n)
		case reflect.Struct:
			sv := v.Field(i)
			if err := decode(p, sv, tag); err != nil {
				return err
			}
		case reflect.Slice, reflect.Array:
		}
	}

	return nil
}
