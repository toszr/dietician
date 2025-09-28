package meal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessIngredients(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {
		ingredients := "Mąka Orkiszowa Jasna, Jaja Kurze"
		expected := []string{
			"Mąka orkiszowa jasna",
			"Jaja kurze",
		}
		assert.Equal(t, expected, ProcessIngredients(ingredients))
	})

	t.Run("with parenthesis", func(t *testing.T) {
		ingredients := "Wanilia (Perły Wanilii (62,5%), Naturalny Koncentrat Waniliowy 37,5%))"
		expected := []string{
			"Wanilia (perły wanilii (62.5%), naturalny koncentrat waniliowy 37.5%))",
		}
		assert.Equal(t, expected, ProcessIngredients(ingredients))
	})

	t.Run("special case 'bez skóry'", func(t *testing.T) {
		ingredients := "Filet Z Piersi Kurczaka, bez skóry"
		expected := []string{
			"Filet z piersi kurczaka (bez skóry)",
		}
		assert.Equal(t, expected, ProcessIngredients(ingredients))
	})
}
