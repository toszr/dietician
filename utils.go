package main

import (
	"encoding/xml"
	"strings"
	"unicode"
)

// getOutputPath returns the output path: if outputPath is empty, replaces inputPath's extension with .md
func getOutputPath(inputPath, outputPath string) string {
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

// Find the first child node with a data-cy attribute (excluding wrappers/ingredients) and return its text
func findFirstDishNameNode(n Node) string {
	if ok, val := getAttr(n, "data-cy"); ok && val == "" {
		text := strings.TrimSpace(n.Content)
		if text != "" {
			return text
		}
	}
	for _, c := range n.Nodes {
		if t := findFirstDishNameNode(c); t != "" {
			return t
		}
	}
	return ""
}

// Find ingredients for a dish node
func findIngredients(n Node) string {
	if _, val := getAttr(n, "data-cy"); val == "IngredientsAndRecipes_span" && strings.TrimSpace(n.Content) != "" {
		return strings.TrimSpace(n.Content)
	}
	for _, c := range n.Nodes {
		if t := findIngredients(c); t != "" {
			return t
		}
	}
	return ""
}

type Node struct {
	XMLName xml.Name
	Attr    []xml.Attr `xml:",any,attr"`
	Nodes   []Node     `xml:",any"`
	Content string     `xml:",chardata"`
}

// getAttr returns (ok, value) where ok is true if the attribute is present
func getAttr(n Node, key string) (bool, string) {
	for _, a := range n.Attr {
		if a.Name.Local == key {
			return true, a.Value
		}
	}
	return false, ""
}

func findMeals(root Node) []Node {
	var meals []Node
	var stack []Node
	stack = append(stack, root)
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, val := getAttr(n, "data-cy"); val == "MealDropdownOptions_div" {
			meals = append(meals, n)
		}
		for i := len(n.Nodes) - 1; i >= 0; i-- {
			stack = append(stack, n.Nodes[i])
		}
	}
	return meals
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

// fixIngredientPercentages joins lines that are broken by percentage numbers (e.g. 62\n5%) -> 62.5%)
func fixIngredientPercentages(ings []string) []string {
	var out []string
	i := 0
	for i < len(ings) {
		ing := strings.TrimSpace(ings[i])
		// If this line ends with a number and the next line starts with a percent, join them
		if i+1 < len(ings) {
			next := strings.TrimSpace(ings[i+1])
			if perc, ok := joinPercent(ing, next); ok {
				out = append(out, perc)
				i += 2
				continue
			}
		}
		out = append(out, ing)
		i++
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
		if buf != "" {
			buf += ", " + ing
		} else {
			buf = ing
		}
		open += strings.Count(ing, "(")
		open -= strings.Count(ing, ")")
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

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
