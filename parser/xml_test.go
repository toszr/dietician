package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseXMLToMarkdown(t *testing.T) {
	t.Run("empty XML", func(t *testing.T) {
		input := ""

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.Error(t, err, "ParseXMLToMarkdown should return an error for empty XML")
		assert.Empty(t, result, "Result should be empty when error occurs")
	})

	t.Run("invalid XML", func(t *testing.T) {
		input := "<invalid xml"

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.Error(t, err, "ParseXMLToMarkdown should return an error for invalid XML")
		assert.Empty(t, result, "Result should be empty when error occurs")
	})

	t.Run("XML with no meals", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?><root><div>No meals here</div></root>`
		expected := ""

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML with no meals")
		assert.Equal(t, expected, result, "Result should be empty string when no meals are found")
	})

	t.Run("XML with single meal and dish", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Śniadanie
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Jajecznica</div>
			<span data-cy="IngredientsAndRecipes_span">jajka 2 szt., masło 10g, sól</span>
		</div>
	</div>
</root>`
		expected := `# Śniadanie

## Jajecznica
**Składniki:**
- Jajka 2 szt.
- Masło 10g
- Sól

`

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML")
		assert.Equal(t, expected, result, "Should parse single meal and dish correctly")
	})

	t.Run("XML with multiple meals and dishes", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Śniadanie
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Owsianka</div>
			<span data-cy="IngredientsAndRecipes_span">płatki owsiane 50g, mleko 200ml</span>
		</div>
	</div>
	<div data-cy="MealDropdownOptions_div">
		Obiad
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Kotlet</div>
			<span data-cy="IngredientsAndRecipes_span">mięso wieprzowe 100g, bułka tarta 20g</span>
		</div>
	</div>
</root>`
		expected := `# Śniadanie

## Owsianka
**Składniki:**
- Płatki owsiane 50g
- Mleko 200ml

# Obiad

## Kotlet
**Składniki:**
- Mięso wieprzowe 100g
- Bułka tarta 20g

`

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML")
		assert.Equal(t, expected, result, "Should parse multiple meals and dishes correctly")
	})

	t.Run("XML with dish without ingredients", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Śniadanie
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Kawa</div>
		</div>
	</div>
</root>`
		expected := `# Śniadanie

## Kawa

`

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML")
		assert.Equal(t, expected, result, "Should handle dishes without ingredients correctly")
	})

	t.Run("XML with uppercase ingredients", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Obiad
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Sałatka</div>
			<span data-cy="IngredientsAndRecipes_span">POMIDOR 100g, OGÓREK 50g</span>
		</div>
	</div>
</root>`
		expected := `# Obiad

## Sałatka
**Składniki:**
- Pomidor 100g
- Ogórek 50g

`

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML")
		assert.Equal(t, expected, result, "Should handle uppercase ingredients correctly")
	})

	t.Run("XML with broken percentage ingredients", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Obiad
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Mięso</div>
			<span data-cy="IngredientsAndRecipes_span">wołowina 62
5%, sól</span>
		</div>
	</div>
</root>`
		expected := `# Obiad

## Mięso
**Składniki:**
- Wołowina 62.5%
- Sól

`

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML")
		assert.Equal(t, expected, result, "Should handle broken percentage ingredients correctly")
	})

	t.Run("XML with 'bez skóry' special case", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Obiad
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Kurczak</div>
			<span data-cy="IngredientsAndRecipes_span">pierś z kurczaka, bez skóry, sól</span>
		</div>
	</div>
</root>`
		expected := `# Obiad

## Kurczak
**Składniki:**
- Pierś z kurczaka (bez skóry)
- Sól

`

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML")
		assert.Equal(t, expected, result, "Should handle 'bez skóry' special case correctly")
	})

	t.Run("XML with meal but no dish name", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Śniadanie
		<div data-cy="dish-tile__wrapper">
			<span data-cy="IngredientsAndRecipes_span">jajka 2 szt.</span>
		</div>
	</div>
</root>`
		expected := `# Śniadanie

`

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML")
		assert.Equal(t, expected, result, "Should handle missing dish names correctly")
	})
}

func TestParseXMLToMarkdownWithComplexIngredients(t *testing.T) {
	t.Run("complex ingredients with parentheses and special cases", func(t *testing.T) {
		input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Obiad
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Kotlet schabowy</div>
			<span data-cy="IngredientsAndRecipes_span">mięso wieprzowe (schab 80,5%), bułka tarta (pszenica), jajko (klasa M, bez skóry), sól morska</span>
		</div>
	</div>
</root>`

		expected := `# Obiad

## Kotlet schabowy
**Składniki:**
- Mięso wieprzowe (schab 80.5%)
- Bułka tarta (pszenica)
- Jajko (klasa m, bez skóry)
- Sól morska

`

		result, err := ParseXMLToMarkdown([]byte(input))

		assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for valid XML")
		assert.Equal(t, expected, result, "Should handle complex ingredients with parentheses and special cases correctly")
	})
}

func TestParseXMLToMarkdownPerformance(t *testing.T) {
	t.Run("large XML with multiple meals and dishes", func(t *testing.T) {
		// Create a larger XML with multiple meals and dishes
		var xmlBuilder strings.Builder
		xmlBuilder.WriteString(`<?xml version="1.0" encoding="UTF-8"?><root>`)

		for i := 0; i < 10; i++ {
			xmlBuilder.WriteString(`<div data-cy="MealDropdownOptions_div">`)
			xmlBuilder.WriteString("Posiłek " + string(rune('A'+i)))

			for j := 0; j < 5; j++ {
				xmlBuilder.WriteString(`<div data-cy="dish-tile__wrapper">`)
				xmlBuilder.WriteString(`<div data-cy="">Danie ` + string(rune('A'+j)) + `</div>`)
				xmlBuilder.WriteString(`<span data-cy="IngredientsAndRecipes_span">składnik1 100g, składnik2 50g, składnik3 25g</span>`)
				xmlBuilder.WriteString(`</div>`)
			}

			xmlBuilder.WriteString(`</div>`)
		}
		xmlBuilder.WriteString(`</root>`)

		input := xmlBuilder.String()

		// Run the function multiple times to test performance
		for i := 0; i < 100; i++ {
			result, err := ParseXMLToMarkdown([]byte(input))
			assert.NoError(t, err, "ParseXMLToMarkdown should not return an error for large valid XML")
			assert.NotEmpty(t, result, "Result should not be empty for large XML with meals")
		}
	})
}

// BenchmarkParseXMLToMarkdown benchmarks the ParseXMLToMarkdown function
func BenchmarkParseXMLToMarkdown(b *testing.B) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<root>
	<div data-cy="MealDropdownOptions_div">
		Śniadanie
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Jajecznica</div>
			<span data-cy="IngredientsAndRecipes_span">jajka 2 szt., masło 10g, sól, pieprz</span>
		</div>
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Tost</div>
			<span data-cy="IngredientsAndRecipes_span">chleb pełnoziarnisty 2 kromki, masło 5g</span>
		</div>
	</div>
	<div data-cy="MealDropdownOptions_div">
		Obiad
		<div data-cy="dish-tile__wrapper">
			<div data-cy="">Kotlet</div>
			<span data-cy="IngredientsAndRecipes_span">mięso wieprzowe 100g, bułka tarta 20g, jajko 1 szt.</span>
		</div>
	</div>
</root>`

	data := []byte(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseXMLToMarkdown(data)
		if err != nil {
			b.Fatalf("ParseXMLToMarkdown() error = %v", err)
		}
	}
}
