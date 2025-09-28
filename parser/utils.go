package parser

import (
	"strings"
	"unicode"
)

// GetOutputPath returns the output path: if outputPath is empty, replaces inputPath's extension with .md
func GetOutputPath(inputPath, outputPath string) string {
	if outputPath != "" {
		return outputPath
	}
	outPath := inputPath
	if dot := strings.LastIndex(outPath, "."); dot != -1 {
		outPath = outPath[:dot] + ".md"
	} else {
		outPath = outPath + ".md"
	}
	return outPath
}

// utf8TitleCase capitalizes only the first letter, rest is lowercased (UTF-8 aware)
func utf8TitleCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// isAllUppercase checks if all letters in the string are uppercase (UTF-8 aware)
func isAllUppercase(s string) bool {
	if len(s) <= 1 {
		return false
	}
	for _, r := range []rune(s) {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	// Check if there's at least one letter
	hasLetter := false
	for _, r := range []rune(s) {
		if unicode.IsLetter(r) {
			hasLetter = true
			break
		}
	}
	return hasLetter
}

// isAllDigits checks if all characters in string are digits
func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
