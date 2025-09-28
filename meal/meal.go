package meal

import "strings"

// Dish represents a single dish with its name and ingredients
type Dish struct {
	Name            string   `json:"dishName"`
	Ingredients     []string `json:"ingredients"`
	IngredientsList string   `json:"ingredientsList"`
}

// Meal represents a meal with its name and dishes
type Meal struct {
	Name   string `json:"mealName"`
	Dishes []Dish `json:"dishes"`
}

// Plan represents the structured data for all meals
type Plan []Meal

// FormatToMarkdown converts a meal Plan to Markdown format
func (p *Plan) FormatToMarkdown() string {
	var sb strings.Builder

	// Iterate through meals in the original order
	for _, meal := range *p {
		sb.WriteString("# " + meal.Name + "\n\n")

		for _, dish := range meal.Dishes {
			sb.WriteString("## " + dish.Name + "\n")
			if len(dish.Ingredients) > 0 {
				sb.WriteString("**Składniki:**\n")
				for _, ing := range dish.Ingredients {
					sb.WriteString("- " + ing + "\n")
				}
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
