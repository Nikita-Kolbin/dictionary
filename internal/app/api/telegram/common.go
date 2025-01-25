package telegram

import "strings"

var needEscapedChars = map[rune]struct{}{
	'_': {}, '*': {}, '[': {}, ']': {}, '(': {}, ')': {}, '~': {}, '`': {}, '>': {},
	'#': {}, '+': {}, '-': {}, '=': {}, '|': {}, '{': {}, '}': {}, '.': {}, '!': {},
}

func escapeFormatChars(msg string) string {
	builder := strings.Builder{}
	runes := []rune(msg)
	for _, r := range runes {
		if _, ok := needEscapedChars[r]; ok {
			builder.WriteRune('\\')
		}
		builder.WriteRune(r)
	}
	return builder.String()
}
