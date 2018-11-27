package main

import "testing"

func PlusOne(num int) int {
	return num + 1
}

func TestPlusOne(t *testing.T) {
	tests := []struct {
		input  int
		output int
	}{
		{1, 2},
		{2, 3},
	}

	for _, tt := range tests {
		actual := PlusOne(tt.input)
		if actual != tt.output {
			t.Fatalf("Failed to PlusOne. actual=%d, expected=%d", actual, tt.output)
		}
	}
}

func Concat(a, b string) string {
	return a + b
}

func TestConcat(t *testing.T) {
	tests := []struct {
		str1   string
		str2   string
		output string
	}{
		{"hoge", "foo", "hogefoo"},
		{"fuga", "-hoge", "fuga-hoge"},
	}

	for _, tt := range tests {
		actual := Concat(tt.str1, tt.str2)
		if actual != tt.output {
			t.Fatalf("Failed to PlusOne. actual=%s, expected=%s", actual, tt.output)
		}
	}
}
