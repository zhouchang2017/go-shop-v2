package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func mapGet(s interface{}, key string) (v interface{}, err error) {
	var (
		i  int64
	)
	//typeOf := reflect.TypeOf(s)
	valueOf := reflect.ValueOf(s)
	switch valueOf.Kind() {
	case reflect.Map:
		value:= valueOf.MapIndex(reflect.ValueOf(key))
		if value.IsValid() {
			v =  valueOf.MapIndex(reflect.ValueOf(key)).Interface()
		}else {
			v = nil
		}
	case reflect.Slice:
		if i, err = strconv.ParseInt(key, 10, 64); err == nil {
			if valueOf.Len() >= int(i) {
				v = valueOf.Index(int(i)).Interface()
			}else{
				err = fmt.Errorf("Index out of bounds. [Index:%d] [Array:%v]", i, s)
			}
		}
	}
	return v, err
}

func MapGet(s interface{}, key string) (v interface{}, err error) {

	attrs := strings.Split(key, ".")
	if attrs[0] == "" {
		return s, nil
	}
	head := attrs[0]
	others := strings.Join(attrs[1:], ".")

	target, err := mapGet(s, head)
	if err == nil {
		return MapGet(target, others)
	}
	return s, err
}
