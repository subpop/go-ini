package ini

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type SyntaxError struct {
	LineNum int
	Line    string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("ini: syntax error: %v: '%v'", e.LineNum, e.Line)
}

// An UnmarshalTypeError describes a INI value that was
// not appropriate for a value of a specific Go type.
type UnmarshalTypeError struct {
	Value  string       // description of INI value - "bool", "array", "number -5"
	Type   reflect.Type // type of Go value it could not be assigned to
	Offset int64        // error occurred after reading Offset bytes
	Struct string       // name of the struct type containing the field
	Field  string       // the full path from root node to the field
}

func (e *UnmarshalTypeError) Error() string {
	if e.Struct != "" || e.Field != "" {
		return "ini: cannot unmarshal " + e.Value + " into Go struct field " + e.Struct + "." + e.Field + " of type " + e.Type.String()
	}
	return "ini: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}

// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
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

type section map[string]interface{}
type ini struct {
	global   section
	sections map[string]section
}

type decodeState struct {
	line         string
	lineNum      int
	inSection    bool
	scanner      *bufio.Scanner
	errorContext struct {
		Struct     reflect.Type
		FieldStack []string
	}
	savedError error
}

func (d *decodeState) init(data []byte) *decodeState {
	d.line = ""
	d.lineNum = 0
	d.scanner = bufio.NewScanner(bytes.NewReader(data))
	d.savedError = nil
	return d
}

// saveError saves the first err it is called with,
// for reporting at the end of the unmarshal.
func (d *decodeState) saveError(err error) {
	if d.savedError == nil {
		d.savedError = d.addErrorContext(err)
	}
}

// addErrorContext returns a new error enhanced with information from d.errorContext
func (d *decodeState) addErrorContext(err error) error {
	if d.errorContext.Struct != nil || len(d.errorContext.FieldStack) > 0 {
		switch err := err.(type) {
		case *UnmarshalTypeError:
			err.Struct = d.errorContext.Struct.Name()
			err.Field = strings.Join(d.errorContext.FieldStack, ".")
			return err
		}
	}
	return err
}

func (d *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	sectionRegexp := regexp.MustCompile(`^\[([\w\s]+)\]$`)
	keyValueRegexp := regexp.MustCompile(`^([\w\s]+)=(.*)$`)

	var content = ini{
		global:   make(section),
		sections: make(map[string]section),
	}
	var current string

	// scan and load values into the content struct
	for d.scanner.Scan() {
		if d.savedError != nil {
			break
		}
		d.line = strings.TrimSpace(d.scanner.Text())

		if d.inSection {
			if sectionRegexp.MatchString(d.line) {
				d.saveError(&SyntaxError{LineNum: d.lineNum, Line: d.line})
				break
			}

			if keyValueRegexp.MatchString(d.line) {
				matches := keyValueRegexp.FindStringSubmatch(d.line)
				if len(matches) == 3 {
					content.sections[current][strings.TrimSpace(matches[1])] = strings.TrimSpace(matches[2])
				}
				continue
			}

			d.inSection = false
			current = ""
		} else {
			if sectionRegexp.MatchString(d.line) {
				// new section
				d.inSection = true
				matches := sectionRegexp.FindStringSubmatch(d.line)
				if len(matches) == 2 {
					content.sections[strings.TrimSpace(matches[1])] = make(section)
					current = matches[1]
					continue
				}
			}

			if keyValueRegexp.MatchString(d.line) {
				// global
				matches := keyValueRegexp.FindStringSubmatch(d.line)
				if len(matches) == 3 {
					content.global[matches[1]] = matches[2]
				}
				continue
			}
		}
	}

	// only structs are supported
	if rv.Elem().Kind() != reflect.Struct {
		d.saveError(&InvalidUnmarshalError{reflect.TypeOf(v)})
	}

	d.value(rv.Elem(), &content, "")

	return d.savedError
}

func (d *decodeState) value(v reflect.Value, content *ini, section string) {
	if v.Kind() == reflect.Ptr {
		v = reflect.ValueOf(v)
	}
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.Type().NumField(); i++ {
			sf := v.Type().Field(i)
			tagSlice := strings.SplitN(sf.Tag.Get("ini"), ",", 2)
			tag := tagSlice[0]

			if sf.Type.Kind() == reflect.Struct {
				f := v.Field(i)
				d.value(f, content, tag)
			} else {
				k := sf.Type.Kind()
				switch k {
				case reflect.String:
					if section == "" {
						v.Field(i).SetString(content.global[tag].(string))
					} else {
						v.Field(i).SetString(content.sections[section][tag].(string))
					}
				case reflect.Int:
					var s string
					if section == "" {
						s = content.global[tag].(string)
					} else {
						s = content.sections[section][tag].(string)
					}
					val, err := strconv.ParseInt(s, 10, 64)

					if err != nil {
						panic(err)
					}
					v.Field(i).SetInt(val)
				default:
					d.saveError(&UnmarshalTypeError{
						Value:  reflect.ValueOf(content.global[tag]).String(),
						Type:   reflect.TypeOf(content.global[tag]),
						Struct: v.Type().Name(),
						Field:  sf.Name,
					})
				}
			}
		}
	}
}

func Unmarshal(data []byte, v interface{}) error {
	var d decodeState
	d.init(data)
	return d.unmarshal(v)
}
