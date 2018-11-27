package main

import "testing"

func PlusOne(num int) int {
	return num + 1
}

func TestMain(t *testing.T) {
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
