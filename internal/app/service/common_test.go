package service

import (
	"testing"
)

func Test_escapeFormatChars(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"abcXYZ", "abcXYZ"},
		{"_", "\\_"},
		{"!*", "\\!\\*"},
		{"abc_def", "abc\\_def"},
		{"test (#1)", "test \\(\\#1\\)"},
		{"[hello] (world)", "\\[hello\\] \\(world\\)"},
		{"Go + Rust = love", "Go \\+ Rust \\= love"},
		{"a|b{c}d", "a\\|b\\{c\\}d"},
		{"multiple...dots!", "multiple\\.\\.\\.dots\\!"},
		{"special_chars_*_[test]", "special\\_chars\\_\\*\\_\\[test\\]"},
	}

	for _, tc := range cases {
		out := escapeFormatChars(tc.input)
		if out != tc.expected {
			t.Errorf("Input: %q | Expected: %q | Got: %q", tc.input, tc.expected, out)
		}
	}
}
