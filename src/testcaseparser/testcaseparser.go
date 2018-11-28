package testcaseparser

import (
	"strings"
)

func Parse(testCaseStr string) []string {
	notIncludeBlank := strings.Replace(testCaseStr, " ", "", -1)
	replacedQuote := strings.Replace(notIncludeBlank, "'", "\"", -1)
	return RemoveOuterBracket(replacedQuote)
}

func RemoveOuterBracket(str string) []string {
	result := []string{}

	foundLeftBracketNum := 0
	wordPos := 0
	for i, c := range str {
		if c == '{' {
			foundLeftBracketNum++

		} else if foundLeftBracketNum == 1 && c == '}' { // if c is last bracket
			result = append(result, str[wordPos+1:i])

		} else if c == '}' {
			foundLeftBracketNum--

		} else if foundLeftBracketNum >= 1 && c == ',' {
			result = append(result, str[wordPos+1:i])
			wordPos = i
		}
	}

	return result
}
