// Package gostruct provides functions to work with structs and their tags.
package gostruct

import (
	"fmt"
	"reflect"
	"regexp"
)

// MapStructFieldTags takes a generic type T and a string key.
// It returns a map of field names to their corresponding tag values for all fields in the struct type T
// that contain a tag key that matches the given key.
// The function uses reflection to get the field tags.
func MapStructFieldTags[T any](key string) map[string]string {
	pt := new(T)
	refValue := reflect.ValueOf(pt).Elem()
	refType := refValue.Type()

	n := refType.NumField()
	res := make(map[string]string, n)

	for i := 0; i < n; i++ {
		fieldType := refType.Field(i)

		if tag, ok := GetTagValue(fieldType.Tag, key); ok && tag != "" {
			res[fieldType.Name] = tag
		}
	}

	return res
}

// GetTagValue is designed because luck of functionality in reflect.Tag.Lookup()
// and help retrieve <value> in given <key> from struct fields
func GetTagValue(tag reflect.StructTag, key string) (string, bool) {
	structTag := fmt.Sprintf("%v", tag)
	tagValue := fmt.Sprintf(`(?s)(?i)\s*(?P<key>%s):\"(?P<value>[^\"]+)\"`, key)

	if match := regexp.MustCompile(tagValue).
		FindStringSubmatch(structTag); match != nil {
		return match[2], true
	}

	return "", false
}
