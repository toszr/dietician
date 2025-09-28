package meal

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ProcessIngredients takes a raw string of ingredients, parses it, and returns a slice of formatted ingredient strings.
func ProcessIngredients(ingredients string) []string {
	fixedIngs := fixIngredientPercentages(strings.Split(ingredients, ","))
	var mergedIngs []string
	for i := 0; i < len(fixedIngs); i++ {
		ing := strings.TrimSpace(fixedIngs[i])
		if (ing == "bez skóry" || ing == "bez skóry)") && len(mergedIngs) > 0 {
			if ing == "bez skóry)" {
				mergedIngs[len(mergedIngs)-1] += ", bez skóry)"
			} else {
				mergedIngs[len(mergedIngs)-1] += " (bez skóry)"
			}
			continue
		}
		ing = smartTitleCase(ing)
		mergedIngs = append(mergedIngs, ing)
	}
	return mergedIngs
}

// smartTitleCase converts a string to title case, but keeps text in parentheses in lowercase.
func smartTitleCase(s string) string {
	if strings.Contains(s, "(") && strings.Contains(s, ")") {
		parts := strings.SplitN(s, "(", 2)
		name := utf8SentenceCase(strings.TrimSpace(parts[0]))
		// The rest of the string, which is inside parenthesis
		rest := cases.Lower(language.Polish).String(parts[1])
		// also fix decimal comma in percentages
		rest = regexp.MustCompile(`(\d+),\s*(\d+)%`).ReplaceAllString(rest, "$1.$2%")
		return name + " (" + rest
	}
	return utf8SentenceCase(s)
}

// fixIngredientPercentages joins lines that are broken by percentage numbers (e.g. 62\n5%) -> 62.5%)
func fixIngredientPercentages(ings []string) []string {
	var out []string
	for _, ing := range ings {
		lines := strings.Split(ing, "\n")
		if len(lines) > 1 {
			var fixedLines []string
			i := 0
			for i < len(lines) {
				line := lines[i]
				if i+1 < len(lines) {
					nextLine := lines[i+1]
					if perc, ok := joinPercent(line, nextLine); ok {
						fixedLines = append(fixedLines, perc)
						i += 2
						continue
					}
				}
				fixedLines = append(fixedLines, line)
				i++
			}
			out = append(out, strings.Join(fixedLines, " "))
		} else {
			out = append(out, ing)
		}
	}

	// Now join lines that are part of the same parenthesis group
	return joinParenthesisGroup(out)
}

// joinParenthesisGroup joins ingredient lines until parentheses are balanced
func joinParenthesisGroup(ings []string) []string {
	var result []string
	var buf string
	open := 0
	for _, ing := range ings {
		trimmedIng := strings.TrimSpace(ing)
		if buf != "" {
			buf += ", " + trimmedIng
		} else {
			buf = trimmedIng
		}
		open += strings.Count(trimmedIng, "(")
		open -= strings.Count(trimmedIng, ")")
		if open <= 0 {
			result = append(result, buf)
			buf = ""
			open = 0
		}
	}
	if buf != "" {
		result = append(result, buf)
	}
	return result
}

// joinPercent joins e.g. "62" and "5%)" into "62.5%)" if next starts with digit(s) and ends with %
func joinPercent(a, b string) (string, bool) {
	// a must end with digits, b must start with digits and end with % or %)
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if len(a) > 0 && len(b) > 0 && isAllDigits(a[len(a)-1:]) && (strings.HasSuffix(b, "%") || strings.HasSuffix(b, "%)")) {
		// Find leading digits in b
		j := 0
		for j < len(b) && b[j] >= '0' && b[j] <= '9' {
			j++
		}
		if j > 0 {
			return a + "." + b[:j] + b[j:], true
		}
	}
	return "", false
}

// isAllUppercase checks if a string is all uppercase, ignoring spaces and parentheses
func isAllUppercase(s string) bool {
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}
	return hasLetter
}

// isAllDigits checks if a string is all digits
func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// utf8SentenceCase converts a string to sentence case, handling UTF-8 characters
func utf8SentenceCase(s string) string {
	if s == "" {
		return ""
	}
	// First, lowercase the entire string.
	lower := cases.Lower(language.Polish).String(s)

	// Capitalize the first letter.
	r, size := utf8.DecodeRuneInString(lower)
	if r == utf8.RuneError {
		return s
	}

	caser := cases.Title(language.Polish)
	return caser.String(string(r)) + lower[size:]
}
