// Package gostruct provides functions to work with structs and their tags.
package gostruct

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
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

// SetField sets the value of a struct to a field.
func SetField(fieldName string, dst any, value any) error {
	fieldVal := reflect.ValueOf(dst).Elem().FieldByName(fieldName)

	src := reflect.ValueOf(value)
	if !src.Type().AssignableTo(fieldVal.Type()) {
		return fmt.Errorf("cannot assign %v to %v", src.Type(), fieldVal.Type())
	}

	fieldVal.Set(src)

	return nil
}

func MakeFromHTMLNode[T any](doc *html.Node, tag string) (T, error) {
	var ent T
	entFields := MapStructFieldTags[T](tag)

	for field, tag := range entFields {
		if tag == "-" || strings.HasSuffix(tag, "[..]") {
			continue
		}

		if node := htmlquery.FindOne(doc, tag); node != nil {
			if val := htmlquery.InnerText(node); val != "" {
				val = strings.TrimSpace(val)

				if strings.Contains(field, "ID") {
					val = path.Base(val)
				}

				if err := SetField(field, &ent, val); err != nil {
					return *new(T), fmt.Errorf("setting field %v: %w", field, err)
				}
			}
		}
	}

	return ent, nil
}

func MakeManyFromHTMLNode[T any](doc *html.Node, sel, tag string) ([]T, error) {
	entFields := MapStructFieldTags[T](tag)

	many := make([]T, 1)

	for _, ancestor := range htmlquery.Find(doc, sel) {
		var one T

		for field, tag := range entFields {
			if tag == "-" || strings.HasSuffix(tag, "[..]") {
				continue
			}

			if node := htmlquery.FindOne(ancestor, tag); node != nil {
				if val := htmlquery.InnerText(node); val != "" {
					val = strings.TrimSpace(val)

					if strings.Contains(field, "ID") {
						val = path.Base(val)
					}

					if err := SetField(field, &one, val); err != nil {
						return nil, fmt.Errorf("setting field %v: %w", field, err)
					}
				}
			}

		}

		many = append(many, one)
	}

	return many, nil
}
