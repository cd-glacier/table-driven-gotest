package testcaseparser

import (
	"strings"
)

func Parse(testCaseStr string) []string {
	escapedSpace := strings.Replace(testCaseStr, "\" \"", "___SPACE___", -1)
	notIncludeBlank := strings.Replace(escapedSpace, " ", "", -1)
	replacedQuote := strings.Replace(notIncludeBlank, "'", "\"", -1)
	revertSpace := strings.Replace(replacedQuote, "___SPACE___", "\" \"", -1)
	return RemoveOuterBracket(revertSpace)
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

		} else if foundLeftBracketNum == 1 && c == ',' {
			result = append(result, str[wordPos+1:i])
			wordPos = i
		}
	}
	return result
}
