package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize[T any](obj T) string {
	return serializeAny(reflect.ValueOf(obj))
}

func serializeAny(val reflect.Value) string {
	if !val.IsValid() {
		return "<nil>"
	}

	switch val.Type().Kind() {
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10)

	case reflect.Float32:
		return strconv.FormatFloat(val.Float(), 'g', -1, 32)

	case reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, 64)

	case reflect.Complex64:
		return strconv.FormatComplex(val.Complex(), 'g', -1, 64)

	case reflect.Complex128:
		return strconv.FormatComplex(val.Complex(), 'g', -1, 128)

	case reflect.Array:
		return serializeArray(val)

	case reflect.Interface, reflect.Pointer:
		if val.IsNil() {
			return "nil"
		}
		return serializeAny(val.Elem())

	case reflect.Map:
		if val.IsNil() {
			return "nil"
		}
		return serializeMap(val)

	case reflect.Slice:
		if val.IsNil() {
			return "nil"
		}
		return serializeArray(val)

	case reflect.String:
		return val.String()

	case reflect.Struct:
		return serializeStruct(val)

	default:
		return val.Type().String()
	}
}

func serializeArray(val reflect.Value) string {
	parts := make([]string, 0, val.Len())
	for i := range val.Len() {
		item := val.Index(i)
		parts = append(parts, serializeAny(item))
	}
	return strings.Join(parts, ",")
}

func serializeMap(val reflect.Value) string {
	parts := make([]string, 0, val.Len())
	iter := val.MapRange()
	for iter.Next() {
		parts = append(parts, serializeAny(iter.Key())+":"+serializeAny(iter.Value()))
	}
	return strings.Join(parts, ",")
}

func serializeStruct(val reflect.Value) string {
	structType := val.Type()
	var sb strings.Builder

	for i := range val.NumField() {
		fieldType := structType.Field(i)
		tag := parsePropertiesTag(fieldType.Tag.Get("properties"))
		if len(tag.name) == 0 {
			tag.name = fieldType.Name
		}

		field := val.Field(i)
		if field.IsValid() && field.IsZero() && tag.omitempty {
			continue
		}

		if sb.Len() > 0 {
			sb.WriteByte('\n')
		}

		sb.WriteString(tag.name)
		sb.WriteByte('=')
		sb.WriteString(serializeAny(field))
	}

	return sb.String()
}

type propertiesTag struct {
	name      string
	omitempty bool
}

func parsePropertiesTag(v string) propertiesTag {
	var res propertiesTag
	for i, v := range strings.Split(v, ",") {
		switch {
		case i == 0:
			res.name = v
		case v == "omitempty":
			res.omitempty = true
		}
	}
	return res
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
