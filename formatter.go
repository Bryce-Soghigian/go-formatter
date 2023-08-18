// Package formatter provides functions to generate Go code that constructs specific values.
package formatter

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

func FormatValue(value reflect.Value) string {
	switch value.Kind() {
	case reflect.String:
		return fmt.Sprintf(`"%s"`, value.String())
	case reflect.Bool:
		return fmt.Sprintf("%t", value.Bool())
	case reflect.Int64:
		return fmt.Sprintf("%d", value.Int())
	case reflect.Slice:
		sliceValues := make([]string, value.Len())
		for i := 0; i < value.Len(); i++ {
			sliceValues[i] = FormatValue(value.Index(i))
		}
		return fmt.Sprintf("[]*%s{%s}", value.Type().Elem(), strings.Join(sliceValues, ", "))
	case reflect.Struct:
		return FormatStruct(value)
	case reflect.Ptr:
		if value.IsNil() {
			return "nil"
		}
		return FormatValue(value.Elem())
	default:
		return fmt.Sprintf("unknown type: %s", value.Type())
	}
}

func FormatStruct(structValue reflect.Value) string {
	var buffer bytes.Buffer

	buffer.WriteString("&")
	buffer.WriteString(structValue.Type().Name())
	buffer.WriteString("{\n")

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structValue.Type().Field(i)

		if field.CanInterface() {
			buffer.WriteString(fmt.Sprintf("%s: %s,\n", fieldType.Name, FormatValue(field)))
		}
	}

	buffer.WriteString("},\n")
	return buffer.String()
}
