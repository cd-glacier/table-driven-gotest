package testcaseparser

import (
	"reflect"
	"testing"
)

func TestRemoveOuterBracket(t *testing.T) {
	tests := []struct {
		input  string
		output []string
	}{
		{"{hoge,{foo}}", []string{"hoge", "{foo}"}},
		{"{[]string{'hoge'},[]string{'foo'}}", []string{"[]string{'hoge'}", "[]string{'foo'}"}},
		{"{[]string{'hoge',' '},[]string{'foo'}}", []string{"[]string{'hoge',' '}", "[]string{'foo'}"}},
	}

	for _, tt := range tests {
		actual := RemoveOuterBracket(tt.input)
		if !reflect.DeepEqual(actual, tt.output) {
			t.Fatalf("Failed to RemoveOuteBracket. actual: %#v, expected: %#v", actual, tt.output)
		}
	}
}
