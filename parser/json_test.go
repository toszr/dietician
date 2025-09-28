package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJSONToMarkdown(t *testing.T) {
	t.Run("valid json", func(t *testing.T) {
		input := `[
			{
				"mealName": "Śniadanie",
				"dishes": [
					{
						"dishName": "Jajecznica",
						"ingredients": ["jajka 2 szt.", "masło 10g", "sól"]
					}
				]
			}
		]`
		expected := "# Śniadanie\n\n## Jajecznica\n**Składniki:**\n- Jajka 2 szt.\n- Masło 10g\n- Sól\n\n"
		result, err := ParseJSONToMarkdown([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("json with ingredientsList", func(t *testing.T) {
		input := `[
			{
				"mealName": "Obiad",
				"dishes": [
					{
						"dishName": "Kurczak w sosie",
						"ingredientsList": "pierś z kurczaka, bez skóry, śmietana 30%"
					}
				]
			}
		]`
		expected := "# Obiad\n\n## Kurczak w sosie\n**Składniki:**\n- Pierś z kurczaka (bez skóry)\n- Śmietana 30%\n\n"
		result, err := ParseJSONToMarkdown([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("invalid json", func(t *testing.T) {
		input := `[{"mealName": "Śniadanie"`
		_, err := ParseJSONToMarkdown([]byte(input))
		assert.Error(t, err)
	})

	t.Run("empty json", func(t *testing.T) {
		input := ``
		_, err := ParseJSONToMarkdown([]byte(input))
		assert.Error(t, err)
	})

	t.Run("json with no meals", func(t *testing.T) {
		input := `[]`
		expected := ""
		result, err := ParseJSONToMarkdown([]byte(input))
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
