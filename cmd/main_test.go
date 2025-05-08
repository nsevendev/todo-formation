package main

import (
	"testing"
)

func TestExtractStringInBacktick(t *testing.T) {
	tests := []struct {
		name	string
		value   string
		expect  string
	}{
		{"test avec des `backtick`", "`test`", "test"},
		{"test sans backtick", "test", ""},
	}

	for _, tt := range tests {
		result := extractStringInBacktick(tt.value)
		if result != tt.expect{
			t.Errorf("%s: got %v, expect %v", tt.name, result, tt.expect)
		}
	}
}