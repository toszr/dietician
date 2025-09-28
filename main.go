package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/toszr/dietician/parser"
)

func main() {
	inputPath := flag.String("input", "", "Path to input XML or JSON file")
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

	// Determine file extension and route to appropriate parser
	ext := strings.ToLower(filepath.Ext(*inputPath))
	var markdown string

	switch ext {
	case ".xml":
		markdown, err = parser.ParseXMLToMarkdown(data)
	case ".json":
		markdown, err = parser.ParseJSONToMarkdown(data)
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported file extension '%s'. Only .xml and .json are supported\n", ext)
		os.Exit(1)
	}

	if err != nil {
		panic(err)
	}

	outPath := parser.GetOutputPath(*inputPath, *outputPath)
	err = os.WriteFile(outPath, []byte(markdown), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("Markdown file generated:", outPath)
}
