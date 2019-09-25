package ini

import (
	"reflect"
	"regexp"
	"strconv"
)

// An UnmarshalTypeError describes a value that was not appropriate for a value
// of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // description of value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Struct string       // name of the struct type containing the field
	Field  string       // name of the field within the struct
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "ini: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}
	return "ini: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// Unmarshal parses the INI-encoded data and stores the result in the value
// pointed to by v. If v is nil or not a pointer to a struct, Unmarshal returns
// an UnmarshalTypeError; INI-encoded data must be encoded into a struct.
//
// Unmarshal uses the inverse of the encodings that Marshal uses, following the
// rules below:
//
// So-called "global" property keys are matched to a struct field within v,
// either by its field name or tag. Values are then decoded according to the
// type of the destination field.
//
// Sections must be unmarshaled into a struct. Unmarshal matches the section
// name to a struct field name or tag. Subsequent property keys are then matched
// against struct field names or tags within the struct.
//
// If a duplicate section name or property key is encountered, Unmarshal will
// allocate a slice according to the number of duplicate keys found, and append
// each value to the slice. If the destination struct field is not a slice type,
// an error is returned.
//
// A struct field tag name may contain a valid regular expression within square
// brackets (i.e. `ini:"[lib.*]"`). If such a tag is detected and the destination
// field is a slice, any section name that matches the regular expression is
// decoded into the destination field as an element in the slice. Under these
// circumstances, the section name is not decoded. However, if a struct field
// named "SectionName" is encountered, the section name decoded into that field.
//
// A struct field tag containing "omitempty" will set the destination field to
// its type's zero value if no corresponding property key was encountered.
func Unmarshal(data []byte, v interface{}) error {
	return unmarshal(data, v, Options{})
}

// UnmarshalWithOptions allows parsing behavior to be configured with an Options
// value.
func UnmarshalWithOptions(data []byte, v interface{}, opts Options) error {
	return unmarshal(data, v, opts)
}

func unmarshal(data []byte, v interface{}, opts Options) error {
	p := newParser(data)
	p.l.opts.allowMultilineEscapeNewline = opts.AllowMultilineValues
	p.l.opts.allowMultilineWhitespacePrefix = opts.AllowMultilineValues
	if err := p.parse(); err != nil {
		return err
	}

	if err := decode(p.ast, reflect.ValueOf(v)); err != nil {
		return err
	}

	return nil
}

// decode sets the underlying values of the value to which rv points to the
// concrete value stored in the corresponding field of ast.
func decode(ast ast, rv reflect.Value) error {
	if rv.Type().Kind() != reflect.Ptr {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(ast).String(),
			Type:  rv.Type(),
		}
	}

	rv = reflect.Indirect(rv)
	if rv.Type().Kind() != reflect.Struct {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(ast).String(),
			Type:  rv.Type(),
		}
	}

	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)

		t := newTag(sf)
		if t.name == "-" {
			continue
		}

		switch sf.Type.Kind() {
		case reflect.Struct:
			sv := rv.Field(i).Addr()
			if t.omitempty && len(ast[t.name]) == 0 {
				continue
			}
			val := ast[t.name][0]
			if err := decodeStruct(val, sv); err != nil {
				return err
			}
		case reflect.Slice:
			sv := rv.Field(i).Addr()
			var val interface{}
			switch sf.Type.Elem().Kind() {
			case reflect.Struct:
				var r *regexp.Regexp
				r, err := t.pattern()
				if err != nil {
					return err
				}
				if r != nil {
					val = make([]section, 0)
					for k, v := range ast {
						if k != "" && r.MatchString(k) {
							for _, s := range v {
								val = append(val.([]section), s)
							}
						}
					}
				} else {
					val = ast[t.name]
				}
			default:
				if t.omitempty && len(ast[""][0].props[t.name].val) == 0 {
					continue
				}
				val = ast[""][0].props[t.name].val
			}
			if err := decodeSlice(val, sv); err != nil {
				return err
			}
		case reflect.String:
			sv := rv.Field(i).Addr()
			if t.omitempty && len(ast[""][0].props[t.name].val) == 0 {
				continue
			}
			val := ast[""][0].props[t.name].val[0]
			if err := decodeString(val, sv); err != nil {
				return err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sv := rv.Field(i).Addr()
			if t.omitempty && len(ast[""][0].props[t.name].val) == 0 {
				continue
			}
			val := ast[""][0].props[t.name].val[0]
			if err := decodeInt(val, sv); err != nil {
				return err
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			sv := rv.Field(i).Addr()
			if t.omitempty && len(ast[""][0].props[t.name].val) == 0 {
				continue
			}
			val := ast[""][0].props[t.name].val[0]
			if err := decodeUint(val, sv); err != nil {
				return err
			}
		case reflect.Float32, reflect.Float64:
			sv := rv.Field(i).Addr()
			if t.omitempty && len(ast[""][0].props[t.name].val) == 0 {
				continue
			}
			val := ast[""][0].props[t.name].val[0]
			if err := decodeFloat(val, sv); err != nil {
				return err
			}
		case reflect.Bool:
			sv := rv.Field(i).Addr()
			if t.omitempty && len(ast[""][0].props[t.name].val) == 0 {
				continue
			}
			val := ast[""][0].props[t.name].val[0]
			if err := decodeBool(val, sv); err != nil {
				return err
			}
		}
	}

	return nil
}

// decodeStruct sets the underlying values of the fields of the value to which
// rv points to the concrete values stored in i. If rv is not a reflect.Ptr,
// decodeStruct returns UnmarshalTypeError.
func decodeStruct(i interface{}, rv reflect.Value) error {
	if reflect.TypeOf(i) != reflect.TypeOf(section{}) || rv.Type().Kind() != reflect.Ptr {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(i).String(),
			Type:  rv.Type(),
		}
	}

	s := i.(section)
	rv = rv.Elem()

	/* magic */
	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)

		t := newTag(sf)
		if t.name == "-" {
			continue
		}

		if t.omitempty && len(s.props[t.name].val) == 0 {
			continue
		}

		switch sf.Type.Kind() {
		case reflect.Slice:
			sv := rv.Field(i).Addr()
			val := s.props[t.name].val
			if err := decodeSlice(val, sv); err != nil {
				return err
			}
		case reflect.String:
			sv := rv.Field(i).Addr()
			var val string
			if sf.Name == "SectionName" {
				val = s.name
			} else {
				val = s.props[t.name].val[0]
			}
			if err := decodeString(val, sv); err != nil {
				return err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sv := rv.Field(i).Addr()
			val := s.props[t.name].val[0]
			if err := decodeInt(val, sv); err != nil {
				return err
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			sv := rv.Field(i).Addr()
			val := s.props[t.name].val[0]
			if err := decodeUint(val, sv); err != nil {
				return err
			}
		case reflect.Float32, reflect.Float64:
			sv := rv.Field(i).Addr()
			val := s.props[t.name].val[0]
			if err := decodeFloat(val, sv); err != nil {
				return err
			}
		case reflect.Bool:
			sv := rv.Field(i).Addr()
			val := s.props[t.name].val[0]
			if err := decodeBool(val, sv); err != nil {
				return err
			}
		}
	}

	return nil
}

// decodeSlice sets the underlying values of the elements of the value to which
// rv points to the concrete values stored in i.
func decodeSlice(v interface{}, rv reflect.Value) error {
	if reflect.TypeOf(v).Kind() != reflect.Slice || rv.Type().Kind() != reflect.Ptr {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(v).String(),
			Type:  rv.Type(),
		}
	}

	rv = rv.Elem()

	var decoderFunc func(interface{}, reflect.Value) error

	switch rv.Type().Elem().Kind() {
	case reflect.String:
		decoderFunc = decodeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		decoderFunc = decodeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		decoderFunc = decodeUint
	case reflect.Struct:
		decoderFunc = decodeStruct
	case reflect.Float32, reflect.Float64:
		decoderFunc = decodeFloat
	case reflect.Bool:
		decoderFunc = decodeBool
	default:
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(v).String(),
			Type:  rv.Type(),
		}
	}

	vv := reflect.MakeSlice(rv.Type(), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())

	for i := 0; i < vv.Len(); i++ {
		sv := vv.Index(i).Addr()
		val := reflect.ValueOf(v).Index(i).Interface()
		if err := decoderFunc(val, sv); err != nil {
			return err
		}
	}

	rv.Set(vv)

	return nil
}

// decodeString sets the underlying value of the value to which rv points to
// the concrete value stored in i. If rv is not a reflect.Ptr, decodeString
// returns UnmarshalTypeError.
func decodeString(i interface{}, rv reflect.Value) error {
	if reflect.TypeOf(i).Kind() != reflect.String || rv.Type().Kind() != reflect.Ptr {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(i).String(),
			Type:  rv.Type(),
		}
	}

	rv.Elem().SetString(i.(string))
	return nil
}

// decodeInt sets the underlying value of the value to which rv points to the
// concrete value stored in i. If rv is not a reflect.Ptr, decodeInt returns
// UnmarshalTypeError.
func decodeInt(i interface{}, rv reflect.Value) error {
	if reflect.TypeOf(i).Kind() != reflect.String || rv.Type().Kind() != reflect.Ptr {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(i).String(),
			Type:  rv.Type(),
		}
	}

	n, err := strconv.ParseInt(i.(string), 10, 64)
	if err != nil {
		switch err.(*strconv.NumError).Err {
		case strconv.ErrRange:
		default:
			return &UnmarshalTypeError{
				Value: reflect.ValueOf(i).String(),
				Type:  rv.Type(),
			}
		}
	}

	rv.Elem().SetInt(n)
	return nil
}

// decodeUint sets the underlying value of the value to which rv points to the
// concrete value stored in i. If rv is not a reflect.Ptr, decodeUint returns
// UnmarshalTypeError.
func decodeUint(i interface{}, rv reflect.Value) error {
	if reflect.TypeOf(i).Kind() != reflect.String || rv.Type().Kind() != reflect.Ptr {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(i).String(),
			Type:  rv.Type(),
		}
	}

	n, err := strconv.ParseUint(i.(string), 10, 64)
	if err != nil {
		switch err.(*strconv.NumError).Err {
		case strconv.ErrRange:
		default:
			return &UnmarshalTypeError{
				Value: reflect.ValueOf(i).String(),
				Type:  rv.Type(),
			}
		}
	}

	rv.Elem().SetUint(n)
	return nil
}

// decodeBool sets the underlying value of the value to which rv points to the
// concrete value stored in i. If rv is not a reflect.Ptr, decodeBool returns
// UnmarshalTypeError.
func decodeBool(i interface{}, rv reflect.Value) error {
	if reflect.TypeOf(i).Kind() != reflect.String || rv.Type().Kind() != reflect.Ptr {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(i).String(),
			Type:  rv.Type(),
		}
	}

	n, err := strconv.ParseBool(i.(string))
	if err != nil {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(i).String(),
			Type:  rv.Type(),
		}
	}

	rv.Elem().SetBool(n)
	return nil
}

// decodeFloat sets the underlying value of the value to which rv points to the
// concrete value stored in i. If rv is not a reflect.Ptr, decodeFloat returns
// UnmarshalTypeError.
func decodeFloat(i interface{}, rv reflect.Value) error {
	if reflect.TypeOf(i).Kind() != reflect.String || rv.Type().Kind() != reflect.Ptr {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(i).String(),
			Type:  rv.Type(),
		}
	}

	n, err := strconv.ParseFloat(i.(string), 64)
	if err != nil {
		return &UnmarshalTypeError{
			Value: reflect.ValueOf(i).String(),
			Type:  rv.Type(),
		}
	}

	rv.Elem().SetFloat(n)
	return nil
}
