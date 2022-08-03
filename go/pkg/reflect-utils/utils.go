// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package reflect_utils

import "reflect"

func GetValueAndType(anyVal interface{}) (value reflect.Value, ttype reflect.Type) {
	value = deref(reflect.ValueOf(anyVal))
	ttype = reflect.TypeOf(anyVal).Elem()
	return value, ttype
}

func deref(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}
	return deref(v.Elem())
}
