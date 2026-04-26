package sdk

import (
	"reflect"
	"strings"
)

// structToMap converts exported struct fields into request parameters.
func structToMap(v any) map[string]any {
	out := map[string]any{}
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return out
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return out
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return out
	}
	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}
		name := f.Name
		if tag := f.Tag.Get("json"); tag != "" {
			name = strings.Split(tag, ",")[0]
		}
		if name == "" || name == "-" {
			continue
		}
		val := rv.Field(i)
		if isZeroValue(val) {
			continue
		}
		out[name] = val.Interface()
	}
	return out
}

// isZeroValue matches the request-builder rule used to omit empty fields.
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isZeroValue(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Map, reflect.Slice, reflect.Pointer, reflect.Interface:
		return v.IsNil()
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	}
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}
