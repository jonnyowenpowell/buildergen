package main

import (
	"unicode"
	"unicode/utf8"
)

// FirstToUpper returns s with the first letter in uppercase.
func FirstToUpper(s string) string {
	r, l := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[l:]
}

// FirstToLower returns s with the first letter in lowercase.
func FirstToLower(s string) string {
	r, l := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[l:]
}
