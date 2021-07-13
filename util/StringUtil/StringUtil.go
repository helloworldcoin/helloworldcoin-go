package StringUtil

import (
	"strconv"
	"unicode/utf8"
)

var BLANKSPACE string = " "

func IsEquals(value1 string, value2 string) bool {
	return value1 == value2
}
func IsNullOrEmpty(value1 string) bool {
	return len(value1) == 0
}
func PrefixPadding(rawValue string, targetLength int, paddingValue string) string {
	target := rawValue
	for utf8.RuneCountInString(target) < targetLength {
		target = paddingValue + target
	}
	return target
}
func Concatenate(value1 string, value2 string) string {
	return value1 + value2
}
func Concatenate3(value1 string, value2 string, value3 string) string {
	return value1 + value2 + value3
}
func ValueOfUint64(number uint64) string {
	return strconv.FormatUint(number, 10)
}
