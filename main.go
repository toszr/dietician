package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	inputPath := flag.String("input", "", "Path to input XML file")
	outputPath := flag.String("output", "", "Path to output Markdown file (optional)")
	flag.Parse()

	if *inputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: -input flag is required")
		os.Exit(1)
	}

	data, err := os.ReadFile(*inputPath)
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		fmt.Fprintln(os.Stderr, "Error: input file is empty")
		os.Exit(1)
	}
	var root Node
	err = xml.Unmarshal(data, &root)
	if err != nil {
		panic(err)
	}
	meals := findMeals(root)
	var sb strings.Builder
	for _, meal := range meals {
		// For meal name, use the first non-empty text node (not data-cy based)
		mealName := ""
		var findFirstTextNode func(Node) string
		findFirstTextNode = func(n Node) string {
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
		mealName = findFirstTextNode(meal)
		if mealName == "" {
			continue
		}
		sb.WriteString("# " + mealName + "\n\n")
		// Find dish wrappers and extract dish names and ingredients
		var stack []Node
		stack = append(stack, meal)
		for len(stack) > 0 {
			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if _, val := getAttr(n, "data-cy"); val == "dish-tile__wrapper" {
				// Dish name: first child node with a data-cy attribute (not wrapper/ingredients)
				dishName := findFirstDishNameNode(n)
				if dishName != "" {
					sb.WriteString("## " + dishName + "\n")
					// Try to find ingredients for this dish
					ingredients := findIngredients(n)
					if ingredients != "" {
						sb.WriteString("**Składniki:**\n")
						// Fix broken percentage numbers in ingredients (e.g. 62\n5%) -> 62.5%)
						fixedIngs := fixIngredientPercentages(strings.Split(ingredients, ","))
						// Special handling: append 'bez skóry' to previous ingredient
						var mergedIngs []string
						for i := 0; i < len(fixedIngs); i++ {
							ing := strings.TrimSpace(fixedIngs[i])
							if ing == "bez skóry" && len(mergedIngs) > 0 {
								mergedIngs[len(mergedIngs)-1] += " (bez skóry)"
								continue
							}
							if ing == strings.ToUpper(ing) && len(ing) > 1 {
								ing = utf8TitleCase(ing)
							}
							mergedIngs = append(mergedIngs, ing)
						}
						for _, ing := range mergedIngs {
							sb.WriteString("- " + ing + "\n")
						}
						sb.WriteString("\n")
					} else {
						sb.WriteString("\n")
					}
				}
			}
			for i := len(n.Nodes) - 1; i >= 0; i-- {
				stack = append(stack, n.Nodes[i])
			}
		}
	}
	outPath := getOutputPath(*inputPath, *outputPath)
	err = os.WriteFile(outPath, []byte(sb.String()), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Markdown file generated:", outPath)
}
