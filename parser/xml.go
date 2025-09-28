package parser

import (
	"encoding/xml"
	"strings"

	"github.com/toszr/dietician/meal"
)

// Node represents an XML node structure
type Node struct {
	XMLName xml.Name
	Attr    []xml.Attr `xml:",any,attr"`
	Nodes   []Node     `xml:",any"`
	Content string     `xml:",chardata"`
}

// ParseXMLToMarkdown parses XML data and returns Markdown output
func ParseXMLToMarkdown(data []byte) (string, error) {
	// Step 1: Parse XML to intermediate structure
	mealPlan, err := parseXMLToMealPlan(data)
	if err != nil {
		return "", err
	}

	// Step 2: Format meal plan to markdown
	return mealPlan.FormatToMarkdown(), nil
}

// parseXMLToMealPlan converts XML data to the structured meal.Plan format
func parseXMLToMealPlan(data []byte) (meal.Plan, error) {
	var root Node
	err := xml.Unmarshal(data, &root)
	if err != nil {
		return meal.Plan{}, err
	}

	var mealPlan meal.Plan
	meals := findMeals(root)

	for _, mealNode := range meals {
		// Get meal name
		mealName := findFirstTextNode(mealNode)
		if mealName == "" {
			continue
		}

		// Parse dishes for this meal
		dishes := parseDishesFromMeal(mealNode)
		// Always add the meal, even if it has no valid dishes
		mealPlan = append(mealPlan, meal.Meal{
			Name:   mealName,
			Dishes: dishes,
		})
	}

	return mealPlan, nil
}

// parseDishesFromMeal extracts all dishes from a meal node
func parseDishesFromMeal(mealNode Node) []meal.Dish {
	var dishes []meal.Dish

	// Find dish wrappers
	var stack []Node
	stack = append(stack, mealNode)
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, val := getAttr(n, "data-cy"); val == "dish-tile__wrapper" {
			dish := parseSingleDish(n)
			if dish.Name != "" {
				dishes = append(dishes, dish)
			}
		}
		for i := len(n.Nodes) - 1; i >= 0; i-- {
			stack = append(stack, n.Nodes[i])
		}
	}

	return dishes
}

// parseSingleDish extracts dish information from a dish node
func parseSingleDish(dishNode Node) meal.Dish {
	dish := meal.Dish{}

	// Get dish name
	dish.Name = findFirstDishNameNode(dishNode)
	if dish.Name == "" {
		return dish
	}

	// Get ingredients
	ingredients := findIngredients(dishNode)
	if ingredients != "" {
		// Process ingredients (fix percentages, handle special cases, etc.)
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
			if isAllUppercase(ing) {
				ing = utf8TitleCase(ing)
			}
			mergedIngs = append(mergedIngs, ing)
		}
		dish.Ingredients = mergedIngs
	}

	return dish
}

func findFirstTextNode(n Node) string {
	if strings.TrimSpace(n.Content) != "" {
		return strings.TrimSpace(n.Content)
	}
	for _, c := range n.Nodes {
		if t := findFirstTextNode(c); t != "" {
			return t
		}
	}
	return ""
}

// findFirstDishNameNode finds the first child node with a data-cy attribute (excluding wrappers/ingredients) and return its text
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

// findIngredients finds ingredients for a dish node
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

// getAttr returns (ok, value) where ok is true if the attribute is present
func getAttr(n Node, key string) (bool, string) {
	for _, a := range n.Attr {
		if a.Name.Local == key {
			return true, a.Value
		}
	}
	return false, ""
}

// findMeals finds all meal nodes in the XML structure
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
