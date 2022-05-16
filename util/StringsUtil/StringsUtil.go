package StringsUtil

import (
	"helloworldcoin-go/util/StringUtil"
	"strings"
)

/*
 @author x.king xdotking@gmail.com
*/

func HasDuplicateElement(values *[]string) bool {
	visited := make(map[string]bool, 0)
	for i := 0; i < len(*values); i++ {
		if visited[(*values)[i]] == true {
			return true
		} else {
			visited[(*values)[i]] = true
		}
	}
	return false
}

func Contains(values *[]string, value string) bool {
	if values != nil && len(*values) != 0 {
		for _, v := range *values {
			if v == value {
				return true
			}
		}
	}
	return false
}

func Split(values string, valueSeparator string) []string {
	if StringUtil.IsEmpty(values) {
		return []string{}
	}
	return strings.Split(values, valueSeparator)
}
