package parser

import (
	"strings"
)

// GetOutputPath returns the output path: if outputPath is empty, replaces inputPath's extension with .md
func GetOutputPath(inputPath, outputPath string) string {
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
