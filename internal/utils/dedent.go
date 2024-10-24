package utils

import (
	"regexp"
	"strings"
)

var (
	rTrailingNewLine = regexp.MustCompile(`^\n`)
	rIndent          = regexp.MustCompile(`^\s+`)
)

func Dedent(str string) (result string) {
	result = rTrailingNewLine.ReplaceAllString(str, "")

	indent := rIndent.FindString(result)
	if indent == "" {
		return result
	}
	rFound := regexp.MustCompile(`(?m)^` + indent)

	return rFound.ReplaceAllString(result, "")
}

func TabToSpaces(str string, count int) string {
	return strings.ReplaceAll(str, "\t", strings.Repeat(" ", count))
}
