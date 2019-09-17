package ini

import (
	"bytes"
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// An UnsupportedTypeError is returned by Marshal when attempting
// to encode an unsupported value type.
type UnsupportedTypeError struct {
	Type reflect.Type
}

func (e *UnsupportedTypeError) Error() string {
	return "ini: unsupported type: " + e.Type.String()
}

// An UnsupportedValueError is returned by Marshal when attempting
// to encode an unsupported value.
type UnsupportedValueError struct {
	Value reflect.Value
}

func (e *UnsupportedValueError) Error() string {
	return "ini: unsupported value: " + e.Value.String()
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
// the encoding.MarshalText, Marshal calls its MarshalText method and encodes the
// result into the INI property. Otherwise Marshal attempts to encode a textual
// representation of the value through string formatting.
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
// encoding if the field as an empty value, defined as false, 0, a nil pointer,
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
// Floating point, integer and Number values enocded as string representations.
//
// String values encode as valid UTF-8 strings.
//
// Pointer values encode as the value pointed to. A nil pointer encodes as an
// empty property value.
//
// Interface values encode as the value contained in the interface. A nil interface
// value encodes as an empty property value.
//
// Channel, complex, map, slice, array and function values cannot be encoded in
// INI. Attempting to encode such a value causes Marshal to return an
// UnsupportedTypeError.
//
// Attempting to marshal any type other than a struct causes Marshal to return an
// UnsupportedValueError.
//
func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer

	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr, reflect.Struct:
		if err := encode(&buf, reflect.ValueOf(v)); err != nil {
			return nil, err
		}
		return bytes.TrimSpace(buf.Bytes()), nil
	default:
		return nil, &UnsupportedTypeError{reflect.ValueOf(v).Type()}
	}
}

func encode(buf *bytes.Buffer, v reflect.Value) error {
	t := v.Type()
	switch v.Kind() {
	case reflect.Ptr:
		return encode(buf, v.Elem())
	case reflect.Struct:
		// First loop over all fields in the struct, skipping struct fields.
		// When inside a struct, this prints out all fields in the struct as
		// properties. Convienently, during the first entry into this function
		// all non-struct fields are written as "global" INI properties.
		for i := 0; i < v.NumField(); i++ {
			sv := v.Field(i)
			if sv.Kind() == reflect.Struct {
				continue
			}
			sf := t.Field(i)
			key, opts := key(sf)
			if key == "-" {
				continue
			}
			if opts["omitempty"] && sv.Interface() == reflect.Zero(sv.Type()).Interface() {
				continue
			}
			buf.WriteString(key)
			buf.WriteRune('=')
			if err := encode(buf, sv); err != nil {
				return err
			}
			buf.WriteRune('\n')
		}
		buf.WriteRune('\n')

		// Next loop over all fields in the struct, skipping all non-struct fields.
		// Each struct field is then encoded recursively through encode().
		for i := 0; i < v.NumField(); i++ {
			sv := v.Field(i)
			if sv.Kind() == reflect.Struct {
				sf := t.Field(i)
				key, opts := key(sf)
				if key == "-" {
					continue
				}
				if opts["omitempty"] && sv.Interface() == reflect.Zero(sv.Type()).Interface() {
					continue
				}
				fmt.Fprintf(buf, "[%v]\n", key)
				if err := encode(buf, sv); err != nil {
					return err
				}
			}
		}
	case reflect.String:
		if err := encodeTextMarshaler(buf, v, func() {
			buf.WriteString(v.String())
		}); err != nil {
			return err
		}
	case reflect.Int:
		if err := encodeTextMarshaler(buf, v, func() {
			buf.WriteString(strconv.FormatInt(v.Int(), 10))
		}); err != nil {
			return err
		}
	case reflect.Bool:
		if err := encodeTextMarshaler(buf, v, func() {
			var s string
			if v.Bool() {
				s = "true"
			} else {
				s = "false"
			}
			buf.WriteString(s)
		}); err != nil {
			return err
		}
	default:
		return &UnsupportedValueError{v}
	}

	return nil
}

func key(sf reflect.StructField) (tag string, tagOpts map[string]bool) {
	tag = sf.Name
	tags := strings.Split(sf.Tag.Get("ini"), ",")
	if tags[0] != "" {
		tag = tags[0]
	}
	tagOpts = make(map[string]bool)
	for _, opt := range tags[1:] {
		tagOpts[opt] = true
	}
	return
}

func encodeTextMarshaler(buf *bytes.Buffer, v reflect.Value, byteWriter func()) error {
	m, ok := v.Interface().(encoding.TextMarshaler)
	if !ok {
		byteWriter()
	} else {
		b, err := m.MarshalText()
		if err != nil {
			return &MarshalerError{v.Type(), err}
		}
		buf.Write(b)
	}

	return nil
}
