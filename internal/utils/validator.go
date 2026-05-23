package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func ValidateStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.New("must validate a struct")
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		
		// If tag `validate:"required"` is present, ensure it's not empty
		tag := fieldType.Tag.Get("validate")
		if tag == "required" {
			if isZeroValue(field) {
				return fmt.Errorf("field %s is required", fieldType.Name)
			}
		}
	}
	return nil
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
