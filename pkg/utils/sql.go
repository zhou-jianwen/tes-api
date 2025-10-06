package utils

import "strings"

// EscapeLikeSpecialChars ...
func EscapeLikeSpecialChars(s string) string {
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
