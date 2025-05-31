package main_test

import (
	"testing"

	"zsb_charm_tmux_urlopen"
)

func TestExtractFirstURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Check this out: http://example.com", "http://example.com"},
		{"Visit https://example.com for more info", "https://example.com"},
		{"No URL here!", ""},
		{"Multiple URLs: http://example.com and https://example.org", "http://example.com"},
		{"Special chars: https://example.com/path?query=1&other=2", "https://example.com/path?query=1&other=2"},
		{"Edge case: just some text", ""},
		{"Another edge case: https://example.com/path#fragment", "https://example.com/path#fragment"},
		{"Another edge case: (https://www.example.cl/vehiculos/autos-veh%C3%ADculo/foo/)", "https://www.example.cl/vehiculos/autos-veh%C3%ADculo/foo/"},
		{"Empty string: ", ""},
		{"Only special chars: !@#$%^&*()", ""},
	}

	for _, test := range tests {
		result, _ := main.ExtractFirstURL(test.input)
		if result != test.expected {
			t.Errorf("For input '%s', expected '%s', but got '%s'", test.input, test.expected, result)
		}
	}
}
