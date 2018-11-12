package utils

import (
	"reflect"
)

// GetStructFields returns array of struct field names
func GetStructFields(v interface{}) []string {
	t := reflect.TypeOf(v)
	numFields := t.NumField()
	names := []string{}
	for i := 0; i < numFields; i++ {
		strField := t.Field(i)
		names = append(names, strField.Name)
	}
	return names
}

// GetTagValue returns tag value for struct field
func GetTagValue(v interface{}, field, tagName string) string {
	t := reflect.TypeOf(v)
	f, _ := t.FieldByName(field)
	return f.Tag.Get(tagName)
}

// IsZeroValue which indicates if value has zero value
func IsZeroValue(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

// GetFieldValue returns value of struct by field name
func GetFieldValue(v interface{}, field string) interface{} {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Interface()
}
