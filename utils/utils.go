package utils

import (
	iohelper "github.com/ndphu/music-downloader/utils/io"
	"strings"
)
// Ngan Test 2
func ArrayContains(arr []int, val int) bool {
	for _, cur := range arr {
		if val == cur {
			return true
		}
	}
	return false
}

func TrimTitle(title string) string {
	return iohelper.CleanupFileName(TrimCDATA(title))
}

func TrimCDATA(input string) string {
	output := input
	for _, ch := range []string{
		"\n",
		" ",
		"\r\n",
	} {
		output = strings.Trim(output, ch)
	}
	return output
}
