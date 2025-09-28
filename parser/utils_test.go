package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOutputPath(t *testing.T) {
	t.Run("with output path", func(t *testing.T) {
		inputPath := "test.xml"
		outputPath := "output.md"
		expected := "output.md"
		assert.Equal(t, expected, GetOutputPath(inputPath, outputPath))
	})

	t.Run("without output path", func(t *testing.T) {
		inputPath := "test.xml"
		outputPath := ""
		expected := "test.md"
		assert.Equal(t, expected, GetOutputPath(inputPath, outputPath))
	})

	t.Run("without output path and no extension", func(t *testing.T) {
		inputPath := "test"
		outputPath := ""
		expected := "test.md"
		assert.Equal(t, expected, GetOutputPath(inputPath, outputPath))
	})

	t.Run("with multiple dots in input path", func(t *testing.T) {
		inputPath := "test.input.xml"
		outputPath := ""
		expected := "test.input.md"
		assert.Equal(t, expected, GetOutputPath(inputPath, outputPath))
	})
}
