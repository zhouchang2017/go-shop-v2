package utils

import (
	"github.com/entropyx/tools/strutils"
	"github.com/jinzhu/inflection"
	"reflect"
	"strings"
	"unicode"
)

func StrToSnake(str string) string {
	return strutils.ToSnakeCase(str)
}

func StrSlug(str string) string {
	bytes := make([]byte, 0, len(str))
	for index, ch := range str {
		if unicode.IsUpper(ch) {
			lower := unicode.ToLower(ch)
			if index == 0 {
				bytes = append(bytes, byte(lower))
			} else {
				bytes = append(bytes, '-')
				bytes = append(bytes, byte(lower))
			}
		} else {
			bytes = append(bytes, byte(ch))
		}
	}
	return string(bytes)
}

func StrPoint(str string) string  {
	bytes := make([]byte, 0, len(str))
	for index, ch := range str {
		if unicode.IsUpper(ch) {
			lower := unicode.ToLower(ch)
			if index == 0 {
				bytes = append(bytes, byte(lower))
			} else {
				bytes = append(bytes, '.')
				bytes = append(bytes, byte(lower))
			}
		} else {
			bytes = append(bytes, byte(ch))
		}
	}
	return string(bytes)
}

func StrToSingular(str string) string {
	return inflection.Singular(str)
}

func StrToPlural(str string) string {
	return inflection.Plural(str)
}

func StrToSlugAndPlural(str string) string {
	return StrSlug(StrToPlural(str))
}

// Brand -> brands , ProductOption -> product_options
func StructNameToSnakeAndPlural(i interface{}) string {
	return StrToSnake(StrToPlural(StructToName(i)))
}

func StructToName(i interface{}) string {
	t := reflect.TypeOf(i)
	split := strings.Split(t.String(), ".")
	name := split[len(split)-1]
	return name
}


func SubString(str string, begin, length int) string {
	rs := []rune(str)
	lth := len(rs)
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length

	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}
