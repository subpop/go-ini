package ini

import (
	"bytes"
	"encoding"
	"reflect"
	"strconv"
)

// A MarshalTypeError represents a type that cannot be encoded in an INI-compatible
// textual format.
type MarshalTypeError struct {
	typ reflect.Type
}

func (e *MarshalTypeError) Error() string {
	return "ini: unsupported type: " + e.typ.String()
}

// A MarshalerError represents an error from calling a MarshalText method.
type MarshalerError struct {
	Type reflect.Type
	Err  error
}

func (e *MarshalerError) Error() string {
	return "ini: error calling MarshalText for type " + e.Type.String() + ": " + e.Err.Error()
}

// Marshal returns the INI encoding of v.
//
// Marshal traverses the value of v recursively. If an encountered value implements
// the encoding.MarshalText interface, Marshal calls its MarshalText method and
// encodes the result into the INI property value. Otherwise Marshal attempts to
// encode a textual representation of the value through string formatting.
//
// The following types are encoded:
//
// Struct values are encoded as INI sections.
// Each exported struct field becomes a property of the section, using the field
// name as the property key, unless the field is omitted for one of the reasons
// given below. The field value encodes as the property value according to the
// rules given below.
//
// The encoding of each struct field can be customized by the format string stored
// under the "ini" key in the struct field's tag.
// The format string gives the name of the field, possibly followed by a
// comma-separated list of options. The name may be empty in order to specify
// options without overriding the default field name.
//
// The "omitempty" option specifies that the field should be omitted from the
// encoding if the field is an empty value, defined as false, 0, a nil pointer,
// a nil interface value, and an empty string.
//
// As a special case, if the field tag is "-", the field is always omitted.
//
// Examples of struct field tags and their meanings:
//
//   // Field appears in INI as key "myName".
//   Field int `ini:"myName"`
//
//   // Field is ignored by this package.
//   Field int `ini:"-"`
//
// Boolean values encode as the string literal "true" or "false".
//
// Floating point, integer and Number values encoded as string representations.
//
// String values encode as valid UTF-8 strings.
//
// Pointer values encode as the value pointed to. A nil pointer encodes as an
// empty property value.
//
// Interface values encode as the value contained in the interface. A nil interface
// value encodes as an empty property value.
//
// Slices and arrays are encoded as a sequential list of properties with
// duplicate keys.
//
// Structs are encoded as a string, the value of which is derived from the
// encoding.TextMarshaler interface. A struct that does not implement this
// interface causes Marshal to return a MarshalTypeError.
//
// Channel, complex, map and function values cannot be encoded in INI.
// Attempting to encode such a value causes Marshal to return a
// MarshalTypeError.
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if err := encode(&buf, reflect.ValueOf(v)); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(buf.Bytes()), nil
}

// encode reflects on the values of rv, encoding them as INI data. If rv is not
// a pointer to a struct, an error is returned. encode makes two passes over
// the struct fields of rv. The first pass skips struct fields that are
// themselve structs, encoding all struct fields as "global" INI properties.
// The second pass then encodes each struct field that *is* a struct as an
// INI section.
func encode(buf *bytes.Buffer, rv reflect.Value) error {
	if rv.Type().Kind() == reflect.Ptr {
		rv = reflect.Indirect(rv)
	}

	if rv.Type().Kind() != reflect.Struct {
		return &MarshalTypeError{typ: rv.Type()}
	}

	// first pass, skipping structs
	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		sv := rv.Field(i)
		t := newTag(sf)

		if t.name == "-" || sf.Type.Kind() == reflect.Struct {
			continue
		}

		if t.omitempty && sv.Interface() == reflect.Zero(sv.Type()).Interface() {
			continue
		}

		if err := encodeProperty(buf, t.name, sv); err != nil {
			return err
		}
	}

	buf.WriteRune('\n')

	// second pass, only structs
	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		sv := rv.Field(i)
		t := newTag(sf)

		if t.name == "-" || sf.Type.Kind() != reflect.Struct {
			continue
		}

		if t.omitempty && sv.Interface() == reflect.Zero(sv.Type()).Interface() {
			continue
		}

		if err := encodeSection(buf, t.name, sv); err != nil {
			return err
		}
	}

	return nil
}

func encodeSection(buf *bytes.Buffer, key string, rv reflect.Value) error {
	if rv.Type().Kind() != reflect.Struct {
		return &MarshalTypeError{typ: rv.Type()}
	}

	buf.WriteString("[" + key + "]\n")

	for i := 0; i < rv.NumField(); i++ {
		sf := rv.Type().Field(i)
		sv := rv.Field(i)
		t := newTag(sf)

		if t.name == "-" {
			continue
		}

		if t.omitempty && sv.Interface() == reflect.Zero(sv.Type()).Interface() {
			continue
		}

		if err := encodeProperty(buf, t.name, sv); err != nil {
			return err
		}
	}

	buf.WriteRune('\n')

	return nil
}

// encodeProperty reflects on the concrete type of rv and encodes it as bytes
// into buf. If rv implements the encoding.TextMarshaler interface, it is used
// to encode the value, otherwise the type is encoded as a string using conversion
// where possible.
func encodeProperty(buf *bytes.Buffer, key string, rv reflect.Value) error {
	var data []byte

	if m, ok := rv.Interface().(encoding.TextMarshaler); ok {
		var err error
		data, err = m.MarshalText()
		if err != nil {
			return &MarshalerError{Type: rv.Type(), Err: err}
		}
	} else {
		switch rv.Type().Kind() {
		case reflect.Slice:
			for i := 0; i < rv.Len(); i++ {
				if err := encodeProperty(buf, key, rv.Index(i)); err != nil {
					return err
				}
			}
		case reflect.Map:
			iter := rv.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()
				subkey := key + "[" + k.String() + "]"
				if err := encodeProperty(buf, subkey, v); err != nil {
					return err
				}
			}
		case reflect.String:
			data = []byte(rv.String())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			data = []byte(strconv.FormatInt(rv.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			data = []byte(strconv.FormatUint(rv.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			data = []byte(strconv.FormatFloat(rv.Float(), 'g', -1, 64))
		case reflect.Bool:
			data = []byte(strconv.FormatBool(rv.Bool()))
		default:
			return &MarshalTypeError{typ: rv.Type()}
		}

	}
	if len(data) > 0 {
		buf.WriteString(key)
		buf.WriteRune('=')
		buf.Write(data)
		buf.WriteRune('\n')
	}
	return nil
}
