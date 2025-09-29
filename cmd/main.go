package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/toszr/dietician/parser"
)

func main() {
	var (
		inputPath  = flag.String("input", "", "Path to the input file (XML or JSON)")
		outputPath = flag.String("output", "", "Path to the output file (Markdown)")
	)
	flag.Parse()

	if *inputPath != "" {
		processFile(*inputPath, *outputPath)
	} else {
		if *outputPath != "" {
			log.Println("Warning: --output flag is ignored when --input is not provided.")
		}
		processSamplesDir("samples")
	}
}

func processFile(inputPath, outputPath string) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("Failed to read input file: %v", err)
	}

	var markdownContent string
	ext := strings.ToLower(filepath.Ext(inputPath))
	switch ext {
	case ".xml":
		markdownContent, err = parser.ParseXMLToMarkdown(data)
	case ".json":
		markdownContent, err = parser.ParseJSONToMarkdown(data)
	default:
		log.Printf("Unsupported file type: %s, skipping", ext)
		return
	}

	if err != nil {
		log.Fatalf("Failed to parse input file '%s': %v", inputPath, err)
	}

	outputFilePath := parser.GetOutputPath(inputPath, outputPath)
	err = os.WriteFile(outputFilePath, []byte(markdownContent), 0644)
	if err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputPath, outputFilePath)
}

func processSamplesDir(inputDir string) {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read samples directory: %v", err)
	}

	mdFiles := make(map[string]bool)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			base := strings.TrimSuffix(file.Name(), ".md")
			mdFiles[base] = true
		}
	}

	processedCount := 0
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == ".xml" || ext == ".json" {
			base := strings.TrimSuffix(file.Name(), ext)
			if !mdFiles[base] {
				inputPath := filepath.Join("samples", file.Name())
				fmt.Printf("Processing %s...\n", inputPath)
				processFile(inputPath, "")
				processedCount++
			}
		}
	}

	if processedCount == 0 {
		fmt.Println("No new files to process in samples/ directory.")
	} else {
		fmt.Printf("Finished processing. %d new markdown file(s) created.\n", processedCount)
	}
}
