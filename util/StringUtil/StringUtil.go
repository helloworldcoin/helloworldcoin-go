package StringUtil

func Concat(value1 string, value2 string) string {
	return value1 + value2
}
func Concat3(value1 string, value2 string, value3 string) string {
	return value1 + value2 + value3
}
func IsNullOrEmpty(value1 string) bool {
	return len(value1) == 0
}
