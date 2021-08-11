package util

import "strings"

var toEscape = [...]int32{'\\', '"', '`', '*', '_', '{', '}', '[', ']', '(', ')', '#', '+', '-', '.', '!'}

// EscapeMarkdown escapes markdown like this: https://meta.stackexchange.com/a/198231
func EscapeMarkdown(markdown string) string {
	var result strings.Builder
	result.Grow(len(markdown))
	for _, c := range markdown {
		if mustBeEscaped(c) {
			result.WriteByte('\\')
		}
		result.WriteRune(c)
	}
	return result.String()
}

func mustBeEscaped(c int32) bool {
	for _, candidate := range toEscape {
		if c == candidate {
			return true
		}
	}
	return false
}
