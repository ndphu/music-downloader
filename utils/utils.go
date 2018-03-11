package utils

import (
	"strings"
)

func ArrayContains(arr []int, val int) bool {
	for _, cur := range arr {
		if val == cur {
			return true
		}
	}
	return false
}

func TrimTitle(title string) string {
	output := title
	for _, ch := range []string{
		"\n",
		" ",
		"\r\n",
	} {
		output = strings.Trim(output, ch)
	}
	return output
}
