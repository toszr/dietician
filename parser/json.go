package parser

import (
	"encoding/json"

	"github.com/toszr/dietician/meal"
)

// ParseJSONToMarkdown parses JSON data and returns Markdown output
func ParseJSONToMarkdown(data []byte) (string, error) {
	var mealPlan meal.Plan
	err := json.Unmarshal(data, &mealPlan)
	if err != nil {
		return "", err
	}

	return mealPlan.FormatToMarkdown(), nil
}
