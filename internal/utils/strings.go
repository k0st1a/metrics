package utils

import "strings"

func AddDotIfNo(s string) string {
	if strings.ContainsRune(s, 46) { // 46 - ascii code of dot
		return s
	}
	return s + "."
}
