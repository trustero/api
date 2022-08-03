package cmd

import (
	"reflect"
	"strings"
)

func expandFieldTag(field reflect.StructField) (tags map[string]string) {
	tags = map[string]string{}
	if val, ok := field.Tag.Lookup("trustero"); ok {
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
