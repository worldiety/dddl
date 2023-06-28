package html

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func UnmarshallValues(dst any, values url.Values) error {
	if reflect.ValueOf(dst).Kind() != reflect.Pointer {
		panic("dst must be a pointer")
	}

	typ := reflect.ValueOf(dst).Elem()

	for key, values := range values {
		field := typ.FieldByName(key)
		if !field.IsValid() {
			return fmt.Errorf("type %T does not have expected form field '%s'", dst, key)
		}

		if err := parseValue(field, values); err != nil {
			return fmt.Errorf("value %v cannot be parsed into field %T.%s: %w", values, dst, key, err)
		}
	}

	return nil
}

func UnmarshallForm(dst any, src *http.Request) error {
	if reflect.ValueOf(dst).Kind() != reflect.Pointer {
		panic("dst must be a pointer")
	}

	typ := reflect.ValueOf(dst).Elem()

	for i := 0; i < typ.NumField(); i++ {
		val := typ.Field(i).Interface()
		if _, ok := val.(*multipart.Form); ok {
			typ.Field(i).Set(reflect.ValueOf(src.MultipartForm))
		}
	}

	for key, values := range src.MultipartForm.Value {
		if strings.HasPrefix(key, "_") {
			continue
		}

		field := typ.FieldByName(key)
		if !field.IsValid() {
			return fmt.Errorf("type %T does not have expected form field '%s'", dst, key)
		}

		if err := parseValue(field, values); err != nil {
			return fmt.Errorf("value %v cannot be parsed into field %T.%s: %w", values, dst, key, err)
		}
	}

	// parse blobs into a byte slice
	for key, headers := range src.MultipartForm.File {
		field := typ.FieldByName(key)
		if !field.IsValid() {
			return fmt.Errorf("type %T does not have expected form file field '%s'", dst, key)
		}

		if len(headers) == 0 {
			continue
		}

		if len(headers) > 1 {
			return fmt.Errorf("cannot parse multiple form file fields '%s' into slice field '%T.%s", key, typ, key)
		}

		r, err := headers[0].Open()
		if err != nil {
			return fmt.Errorf("cannot read file '%s': %w", key, err)
		}

		buf, err := io.ReadAll(r)
		if err != nil {
			return fmt.Errorf("cannot read full file '%s': %w", key, err)
		}

		field.SetBytes(buf)
	}
	return nil
}

func parseValue(field reflect.Value, value []string) error {
	switch field.Type().Kind() {
	case reflect.String:
		if len(value) != 1 {
			return fmt.Errorf("string field has no corresponding value")
		}

		field.SetString(value[0])
	case reflect.Float32, reflect.Float64:
		if len(value) != 1 {
			return fmt.Errorf("float64 field has no corresponding value")
		}

		v, err := strconv.ParseFloat(value[0], 64)
		if err != nil {
			return err
		}
		field.SetFloat(v)
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		if len(value) != 1 {
			return fmt.Errorf("int64 field has no corresponding value")
		}

		v, err := strconv.ParseInt(value[0], 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(v)
	case reflect.Bool:
		if len(value) != 1 {
			return fmt.Errorf("bool field has no corresponding value")
		}

		v, err := strconv.ParseBool(value[0])
		if err != nil {
			return err
		}
		field.SetBool(v)
	case reflect.Slice:
		switch field.Type().Elem().Kind() {
		case reflect.String:
			field.Set(reflect.ValueOf(value))
		default:
			return fmt.Errorf("unsupported slice field type: %v", field)
		}
	default:
		return fmt.Errorf("unsupported field type: %v", field)
	}

	return nil
}
