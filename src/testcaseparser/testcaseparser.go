package testcaseparser

import (
	"strings"
)

func Parse(testCaseStr string) []string {
	notIncludeBlank := strings.Replace(testCaseStr, " ", "", -1)
	return RemoveOuterBracket(notIncludeBlank)
}

func RemoveOuterBracket(str string) []string {
	result := []string{}

	foundLeftBracketNum := 0
	wordPos := 0
	for i, c := range str {
		if c == '{' {
			wordPos = i + 1
			foundLeftBracketNum++
		} else if foundLeftBracketNum >= 2 && c == '}' {
			foundLeftBracketNum--
			result = append(result, str[wordPos-1:i+1])
		}

		if foundLeftBracketNum >= 0 && c == ',' {
			result = append(result, str[wordPos:i])
			wordPos = i
		}
	}

	return result
}
