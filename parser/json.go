package parser

import (
	"encoding/json"
	"strings"
)

// ParseJSONToMarkdown parses JSON data and returns Markdown output
// This function has the same API as ParseXMLToMarkdown for consistency
func ParseJSONToMarkdown(data []byte) (string, error) {
	// Parse the JSON data
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return "", err
	}

	var sb strings.Builder

	// TODO: Implement JSON parsing logic based on the JSON structure
	// For now, this is a placeholder that returns a basic message
	sb.WriteString("# JSON Parsing\n\n")
	sb.WriteString("JSON parsing is not yet implemented. The JSON structure needs to be analyzed first.\n\n")
	sb.WriteString("Raw JSON data preview:\n")

	// Pretty print first 500 characters of JSON for debugging
	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err == nil {
		preview := string(prettyJSON)
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		sb.WriteString("```json\n")
		sb.WriteString(preview)
		sb.WriteString("\n```\n")
	}

	return sb.String(), nil
}
