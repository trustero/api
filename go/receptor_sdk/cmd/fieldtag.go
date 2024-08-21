// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	"reflect"
	"strings"
)

const (
	tagName          = "trustero"
	idField          = "id"
	displayField     = "display"
	orderField       = "order"
	placeholderField = "placeholder"
	controlTestField = "check"
	methodField      = "method"
	inputTypeField   = "input_type"
)

func expandFieldTag(field reflect.StructField) (tags map[string]string) {
	tags = map[string]string{}
	if val, ok := field.Tag.Lookup(tagName); ok {
		tags = getTags(val)
	}
	return
}

func getTags(str string) (m map[string]string) {
	m = map[string]string{}

	if len(str) == 0 {
		return
	}

	if !strings.Contains(str, ";") {
		k, v := getKVPair(str)
		m[k] = v
		return
	}

	pairs := strings.Split(str, ";")
	for _, pair := range pairs {
		if len(pair) > 0 {
			k, v := getKVPair(pair)
			m[k] = v
		}
	}
	return
}

func getKVPair(str string) (k, v string) {
	if strings.Contains(str, ":") {
		kv := strings.Split(str, ":")
		k = kv[0]
		v = kv[1]
	} else {
		k = str
		v = ""
	}
	return
}

func getTagField(tags map[string]string, fieldName, fieldDefault string) string {
	if display, ok := tags[fieldName]; ok {
		return display
	}
	return fieldDefault
}
