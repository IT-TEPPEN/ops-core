package parser

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// FrontmatterData represents the parsed frontmatter metadata from a Markdown document.
type FrontmatterData struct {
	Title       string                   `yaml:"title"`
	Owner       string                   `yaml:"owner"`
	Type        string                   `yaml:"type"`
	Tags        []string                 `yaml:"tags"`
	Variables   []VariableDefinitionYAML `yaml:"variables,omitempty"`
	Content     string                   // The markdown content without frontmatter
	RawFrontmatter string                // The raw YAML frontmatter string
}

// VariableDefinitionYAML represents a variable definition in YAML frontmatter.
type VariableDefinitionYAML struct {
	Name         string      `yaml:"name"`
	Label        string      `yaml:"label"`
	Description  string      `yaml:"description,omitempty"`
	Type         string      `yaml:"type"`
	Required     bool        `yaml:"required"`
	DefaultValue interface{} `yaml:"default_value,omitempty"`
}

// FrontmatterParser defines the interface for parsing frontmatter from Markdown documents.
type FrontmatterParser interface {
	Parse(markdown string) (*FrontmatterData, error)
}

// frontmatterParser implements the FrontmatterParser interface.
type frontmatterParser struct{}

// NewFrontmatterParser creates a new FrontmatterParser.
func NewFrontmatterParser() FrontmatterParser {
	return &frontmatterParser{}
}

// Parse extracts and parses frontmatter from a Markdown document.
// It expects frontmatter to be delimited by "---" at the beginning of the document.
func (p *frontmatterParser) Parse(markdown string) (*FrontmatterData, error) {
	if markdown == "" {
		return nil, fmt.Errorf("markdown content is empty")
	}

	// Check if the document starts with frontmatter delimiter
	if !strings.HasPrefix(markdown, "---\n") && !strings.HasPrefix(markdown, "---\r\n") {
		return nil, fmt.Errorf("frontmatter not found: document must start with '---'")
	}

	// Remove the first "---" delimiter
	markdown = strings.TrimPrefix(markdown, "---\n")
	markdown = strings.TrimPrefix(markdown, "---\r\n")

	// Find the closing "---" delimiter
	closingDelimiterIndex := strings.Index(markdown, "\n---\n")
	if closingDelimiterIndex == -1 {
		closingDelimiterIndex = strings.Index(markdown, "\n---\r\n")
	}
	if closingDelimiterIndex == -1 {
		closingDelimiterIndex = strings.Index(markdown, "\r\n---\r\n")
	}

	if closingDelimiterIndex == -1 {
		return nil, fmt.Errorf("frontmatter closing delimiter '---' not found")
	}

	// Extract frontmatter and content
	rawFrontmatter := markdown[:closingDelimiterIndex]
	content := strings.TrimLeft(markdown[closingDelimiterIndex+5:], "\n\r") // +5 to skip "\n---\n"

	// Parse YAML frontmatter
	var data FrontmatterData
	if err := yaml.Unmarshal([]byte(rawFrontmatter), &data); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter YAML: %w", err)
	}

	// Set the content and raw frontmatter
	data.Content = content
	data.RawFrontmatter = rawFrontmatter

	// Validate required fields
	if err := p.validate(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// validate checks that required frontmatter fields are present and valid.
func (p *frontmatterParser) validate(data *FrontmatterData) error {
	var errors []string

	if strings.TrimSpace(data.Title) == "" {
		errors = append(errors, "title is required")
	}

	if strings.TrimSpace(data.Owner) == "" {
		errors = append(errors, "owner is required")
	}

	if strings.TrimSpace(data.Type) == "" {
		errors = append(errors, "type is required")
	} else {
		// Validate type is either "procedure" or "knowledge"
		if data.Type != "procedure" && data.Type != "knowledge" {
			errors = append(errors, fmt.Sprintf("type must be 'procedure' or 'knowledge', got '%s'", data.Type))
		}
	}

	// Validate variables if present and type is procedure
	if data.Type == "procedure" && len(data.Variables) > 0 {
		for i, v := range data.Variables {
			if strings.TrimSpace(v.Name) == "" {
				errors = append(errors, fmt.Sprintf("variable[%d].name is required", i))
			}
			if strings.TrimSpace(v.Label) == "" {
				errors = append(errors, fmt.Sprintf("variable[%d].label is required", i))
			}
			if strings.TrimSpace(v.Type) == "" {
				errors = append(errors, fmt.Sprintf("variable[%d].type is required", i))
			} else {
				// Validate variable type
				validTypes := map[string]bool{"string": true, "number": true, "boolean": true, "date": true}
				if !validTypes[v.Type] {
					errors = append(errors, fmt.Sprintf("variable[%d].type must be one of [string, number, boolean, date], got '%s'", i, v.Type))
				}
			}
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("frontmatter validation failed: %s", strings.Join(errors, "; "))
	}

	return nil
}
