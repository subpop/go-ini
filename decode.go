package ini

import (
	"reflect"
	"strconv"
)

// An UnmarshalTypeError describes a value that was not appropriate for a value
// of a specific Go type.
type UnmarshalTypeError struct {
	val string       // description of value - "bool", "array", "number -5"
	typ reflect.Type // type of Go value it could not be assigned to
	str string       // name of the struct type containing the field
	fld string       // name of the field within the struct
}

func (e *UnmarshalTypeError) Error() string {
	if e.str != "" || e.fld != "" {
		return "ini: cannot unmarshal " + e.val + " into Go struct field " + e.str + "." + e.fld + " of type " + e.typ.String()
	}
	return "ini: cannot unmarshal " + e.val + " into Go value of type " + e.typ.String()
}

// A DecodeError describes an error that was encountered while parsing a value
// to a specific Go type.
type DecodeError struct {
	err error // underlying parse error
}

func (e DecodeError) Error() string {
	return e.err.Error()
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
// A struct field tag name may be a single asterisk (colloquially known as the
// "wildcard" character). If such a tag is detected and the destination
// field is a slice of structs, all sections are decoded into the destination
// field as an element in the slice. If a struct field named "ININame" is
// encountered, the section name decoded into that field.
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
	p.l.opts.allowNumberSignComments = opts.AllowNumberSignComments
	if err := p.parse(); err != nil {
		return err
	}

	return decode(p.tree, reflect.ValueOf(v))
}

// decode sets the underlying values of the fields of the value to which rv
// points to the parsed values stored in the corresponding field of tree. It
// panics if rv is not a reflect.Ptr to a struct.
func decode(tree parseTree, rv reflect.Value) error {
	rv = rv.Elem()

	// Decode global properties first. By treating rv as the struct to decode
	// into, we ignore any struct fields that are structs.
	if err := decodeStruct(tree.global, rv.Addr()); err != nil {
		return err
	}

	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		sv := rv.Field(i).Addr()

		t := newTag(sf)
		if t.name == "-" {
			continue
		}

		sections, err := tree.get(t.name)
		if err != nil {
			return err
		}

		switch sf.Type.Kind() {
		case reflect.Struct:
			if err := decodeStruct(sections[0], sv); err != nil {
				return err
			}
		case reflect.Slice:
			if sf.Type.Elem().Kind() == reflect.Struct {
				if err := decodeSliceStruct(sections, sv); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// decodeStruct sets the underlying values of the fields of the value to which
// rv points to the parsed values of s. It panics if rv is not a reflect.Ptr to
// a struct.
func decodeStruct(s section, rv reflect.Value) error {
	rv = rv.Elem()

	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		sv := rv.Field(i).Addr()

		t := newTag(sf)
		if t.name == "-" {
			continue
		}

		var val interface{}

		prop, err := s.get(t.name)
		if err != nil {
			return err
		}
		vals := prop.get("")

		var decoderFunc func(string, reflect.Value) error

		switch sf.Type.Kind() {
		case reflect.Slice:
			if sf.Type.Elem().Kind() != reflect.Struct {
				if err := decodeSlice(vals, sv); err != nil {
					return err
				}
			}
			continue
		case reflect.Map:
			if err := decodeMap(*prop, sv); err != nil {
				return err
			}
			continue
		case reflect.String:
			decoderFunc = decodeString
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			decoderFunc = decodeInt
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			decoderFunc = decodeUint
		case reflect.Float32, reflect.Float64:
			decoderFunc = decodeFloat
		case reflect.Bool:
			decoderFunc = decodeBool
		}

		if sf.Name == "ININame" {
			vals = append(vals, s.name)
		}
		if len(vals) == 0 {
			continue
		}
		val = vals[0]

		if err := decoderFunc(val.(string), sv); err != nil {
			return err
		}
	}

	return nil
}

// decodeSliceStruct sets the underlying values of the fields of the elements to
// which rv points to the parsed values of s. It pancis if rv is not a
// reflect.Ptr to a slice of structs.
func decodeSliceStruct(s []section, rv reflect.Value) error {
	rv = rv.Elem()

	vv := reflect.MakeSlice(rv.Type(), len(s), cap(s))

	for i := 0; i < vv.Len(); i++ {
		sv := vv.Index(i).Addr()
		if err := decodeStruct(s[i], sv); err != nil {
			return err
		}
	}

	rv.Set(vv)

	return nil
}

// decodeSlice sets the underlying values of the elements of the value to which
// rv points to the parsed values of s. It panics if rv is not a reflect.Ptr to
// a slice.
func decodeSlice(s []string, rv reflect.Value) error {
	rv = rv.Elem()

	var decoderFunc func(string, reflect.Value) error

	switch rv.Type().Elem().Kind() {
	case reflect.String:
		decoderFunc = decodeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		decoderFunc = decodeInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		decoderFunc = decodeUint
	case reflect.Float32, reflect.Float64:
		decoderFunc = decodeFloat
	case reflect.Bool:
		decoderFunc = decodeBool
	default:
		return &UnmarshalTypeError{
			val: reflect.ValueOf(s).String(),
			typ: rv.Type(),
		}
	}

	vv := reflect.MakeSlice(rv.Type(), len(s), cap(s))

	for i := 0; i < vv.Len(); i++ {
		sv := vv.Index(i).Addr()
		if err := decoderFunc(s[i], sv); err != nil {
			return err
		}
	}

	rv.Set(vv)

	return nil
}

// decodeMap sets the underlying keys and values of the elements of the value to
// which rv points to the parsed values of p. It panics if rv is not a
// reflect.Ptr to a map[string]interface{}.
func decodeMap(p property, rv reflect.Value) error {
	rv = rv.Elem()

	vv := reflect.MakeMap(rv.Type())

	for k, v := range p.vals {
		mv := reflect.New(rv.Type().Elem())

		var decoderFunc func(string, reflect.Value) error

		switch rv.Type().Elem().Kind() {
		case reflect.Slice:
			if err := decodeSlice(v, mv); err != nil {
				return err
			}
			vv.SetMapIndex(reflect.ValueOf(k), mv.Elem())
			continue
		case reflect.String:
			decoderFunc = decodeString
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			decoderFunc = decodeInt
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			decoderFunc = decodeUint
		case reflect.Float32, reflect.Float64:
			decoderFunc = decodeFloat
		case reflect.Bool:
			decoderFunc = decodeBool
		default:
			return &UnmarshalTypeError{
				val: reflect.ValueOf(p).String(),
				typ: rv.Type(),
			}
		}

		if len(v) == 0 {
			continue
		}

		if err := decoderFunc(v[0], mv); err != nil {
			return err
		}
		vv.SetMapIndex(reflect.ValueOf(k), mv.Elem())
	}

	rv.Set(vv)

	return nil
}

// decodeString sets the underlying value of the value to which rv points to
// the parsed value of s. It panics if rv is not a reflect.Ptr to a string.
func decodeString(s string, rv reflect.Value) error {
	rv.Elem().SetString(s)
	return nil
}

// decodeInt sets the underlying value of the value to which rv points to the
// parsed value of s. It panics if rv is not a reflect.Ptr to a int64.
func decodeInt(s string, rv reflect.Value) error {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return &DecodeError{err}
	}

	rv.Elem().SetInt(n)
	return nil
}

// decodeUint sets the underlying value of the value to which rv points to the
// parsed value of s. It panics if rv is not a reflect.Ptr to a uint.
func decodeUint(s string, rv reflect.Value) error {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return &DecodeError{err}
	}

	rv.Elem().SetUint(n)
	return nil
}

// decodeBool sets the underlying value of the value to which rv points to the
// parsed value of s. It panics if rv is not a reflect.Ptr to a bool.
func decodeBool(s string, rv reflect.Value) error {
	n, err := strconv.ParseBool(s)
	if err != nil {
		return &DecodeError{err}
	}

	rv.Elem().SetBool(n)
	return nil
}

// decodeFloat sets the underlying value of the value to which rv points to the
// parsed value of s. It panics if rv is not a reflect.Ptr to a float64.
func decodeFloat(s string, rv reflect.Value) error {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return &DecodeError{err}
	}

	rv.Elem().SetFloat(n)
	return nil
}
