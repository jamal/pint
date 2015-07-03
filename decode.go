package pint

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

type field struct {
	name    string
	typ     reflect.Type
	index   []int
	options tagOptions
}

var fieldCache struct {
	sync.RWMutex
	m map[reflect.Type][]field
}

func getTypeOf(v interface{}) (reflect.Type, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return t, fmt.Errorf("pint: trying to parse a non-struct type %v", t.Kind())
	}
	return t, nil
}

func getValueOf(v interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return val, fmt.Errorf("pint: trying to parse a non-struct type %v", val.Kind())
	}
	return val, nil
}

func getFields(t reflect.Type) ([]field, error) {
	fieldCache.RLock()
	fields := fieldCache.m[t]
	fieldCache.RUnlock()
	if fields != nil {
		return fields, nil
	}

	fieldCache.Lock()
	if fieldCache.m == nil {
		fieldCache.m = make(map[reflect.Type][]field)
	}
	fieldCache.Unlock()

	fields = make([]field, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		sfield := t.Field(i)
		field := field{}
		tag, options := parseTag(sfield.Tag.Get("pint"))
		if tag != "" {
			field.name = tag
		} else {
			field.name = sfield.Name
		}
		field.typ = sfield.Type
		field.index = sfield.Index
		field.options = options
		fields = append(fields, field)
	}

	fieldCache.Lock()
	fieldCache.m[t] = fields
	fieldCache.Unlock()

	return fields, nil
}

// Parse will parse the http.Request form and look for v's fields. It uses
// http.FormValue internally. It will return an error if an invalid v is passed
// and an *ErrValidate if a validation error occurs. Use ErrValidate.String()
// when returning a validation error to the client.
func Parse(r *http.Request, v interface{}) error {
	t, err := getTypeOf(v)
	if err != nil {
		return err
	}
	fields, err := getFields(t)
	if err != nil {
		return err
	}

	// We get the value to write to the struct
	val, err := getValueOf(v)
	if err != nil {
		return err
	}
	for _, field := range fields {
		formval := r.FormValue(field.name)
		// Check if the field is required
		if formval == "" && !field.options.contains("omitempty") {
			return &ErrValidate{fmt.Sprintf("%s cannot be empty", field.name)}
		}

		// Check if we should apply a custom validate/formatter to this field
		formval, err = formatField(formval, field)
		if err != nil {
			return err
		}

		// Parse the value type, run through any type-validators and set the
		// struct value
		switch field.typ.Kind() {
		case reflect.String:
			val.FieldByIndex(field.index).SetString(formval)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intval, err := strconv.ParseInt(formval, 10, 64)
			if err != nil {
				return err
			}
			if err = validateInt(intval, field); err != nil {
				return err
			}
			val.FieldByIndex(field.index).SetInt(intval)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uintval, err := strconv.ParseUint(formval, 10, 64)
			if err != nil {
				return err
			}
			if err = validateInt(int64(uintval), field); err != nil {
				return err
			}
			val.FieldByIndex(field.index).SetUint(uintval)
		case reflect.Float32, reflect.Float64:
			floatval, err := strconv.ParseFloat(formval, 64)
			if err != nil {
				return err
			}
			if err = validateInt(int64(floatval), field); err != nil {
				return err
			}
			val.FieldByIndex(field.index).SetFloat(floatval)
		case reflect.Bool:
			boolval := formval == "true" || formval == "1"
			val.FieldByIndex(field.index).SetBool(boolval)
		}
	}
	return nil
}

func formatField(val string, field field) (string, error) {
	if handlerName, ok := field.options.get("format"); ok {
		// TODO: Should we return an error if the validate handler is not found?
		if handler, ok := formatHandlers[handlerName]; ok {
			return handler(val)
		}
	}
	return val, nil
}
