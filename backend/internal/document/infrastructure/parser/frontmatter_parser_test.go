package parser

import (
	"testing"
)

func TestFrontmatterParser_Parse(t *testing.T) {
	parser := NewFrontmatterParser()

	tests := []struct {
		name        string
		markdown    string
		wantErr     bool
		errContains string
		validate    func(t *testing.T, data *FrontmatterData)
	}{
		{
			name: "valid procedure document with variables",
			markdown: `---
title: Database Backup Procedure
owner: database-team
type: procedure
tags:
  - database
  - backup
variables:
  - name: server_name
    label: Server Name
    description: The target server name
    type: string
    required: true
  - name: port
    label: Port Number
    type: number
    required: false
    default_value: 5432
---

# Database Backup Procedure

This is the content of the document.
`,
			wantErr: false,
			validate: func(t *testing.T, data *FrontmatterData) {
				if data.Title != "Database Backup Procedure" {
					t.Errorf("expected title 'Database Backup Procedure', got '%s'", data.Title)
				}
				if data.Owner != "database-team" {
					t.Errorf("expected owner 'database-team', got '%s'", data.Owner)
				}
				if data.Type != "procedure" {
					t.Errorf("expected type 'procedure', got '%s'", data.Type)
				}
				if len(data.Tags) != 2 {
					t.Errorf("expected 2 tags, got %d", len(data.Tags))
				}
				if len(data.Variables) != 2 {
					t.Errorf("expected 2 variables, got %d", len(data.Variables))
				}
				if data.Variables[0].Name != "server_name" {
					t.Errorf("expected first variable name 'server_name', got '%s'", data.Variables[0].Name)
				}
				if data.Content != "# Database Backup Procedure\n\nThis is the content of the document.\n" {
					t.Errorf("content mismatch: got '%s'", data.Content)
				}
			},
		},
		{
			name: "valid knowledge document without variables",
			markdown: `---
title: System Architecture
owner: engineering
type: knowledge
tags:
  - architecture
  - documentation
---

# System Architecture

Architecture overview...
`,
			wantErr: false,
			validate: func(t *testing.T, data *FrontmatterData) {
				if data.Title != "System Architecture" {
					t.Errorf("expected title 'System Architecture', got '%s'", data.Title)
				}
				if data.Type != "knowledge" {
					t.Errorf("expected type 'knowledge', got '%s'", data.Type)
				}
				if len(data.Variables) != 0 {
					t.Errorf("expected 0 variables for knowledge type, got %d", len(data.Variables))
				}
			},
		},
		{
			name:        "missing frontmatter",
			markdown:    "# Just a heading\n\nNo frontmatter here.",
			wantErr:     true,
			errContains: "frontmatter not found",
		},
		{
			name: "missing closing delimiter",
			markdown: `---
title: Test
owner: team
type: procedure

# Content without closing delimiter
`,
			wantErr:     true,
			errContains: "frontmatter closing delimiter '---' not found",
		},
		{
			name: "missing required field - title",
			markdown: `---
owner: team
type: procedure
tags: []
---

# Content
`,
			wantErr:     true,
			errContains: "title is required",
		},
		{
			name: "missing required field - owner",
			markdown: `---
title: Test Document
type: procedure
tags: []
---

# Content
`,
			wantErr:     true,
			errContains: "owner is required",
		},
		{
			name: "missing required field - type",
			markdown: `---
title: Test Document
owner: team
tags: []
---

# Content
`,
			wantErr:     true,
			errContains: "type is required",
		},
		{
			name: "invalid type value",
			markdown: `---
title: Test Document
owner: team
type: invalid_type
tags: []
---

# Content
`,
			wantErr:     true,
			errContains: "type must be 'procedure' or 'knowledge'",
		},
		{
			name: "invalid variable - missing name",
			markdown: `---
title: Test Document
owner: team
type: procedure
tags: []
variables:
  - label: Server
    type: string
    required: true
---

# Content
`,
			wantErr:     true,
			errContains: "variable[0].name is required",
		},
		{
			name: "invalid variable - missing label",
			markdown: `---
title: Test Document
owner: team
type: procedure
tags: []
variables:
  - name: server
    type: string
    required: true
---

# Content
`,
			wantErr:     true,
			errContains: "variable[0].label is required",
		},
		{
			name: "invalid variable - invalid type",
			markdown: `---
title: Test Document
owner: team
type: procedure
tags: []
variables:
  - name: server
    label: Server
    type: invalid_type
    required: true
---

# Content
`,
			wantErr:     true,
			errContains: "variable[0].type must be one of",
		},
		{
			name:        "empty markdown",
			markdown:    "",
			wantErr:     true,
			errContains: "markdown content is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := parser.Parse(tt.markdown)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing '%s', got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing '%s', got '%s'", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, data)
			}
		})
	}
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

// indexOf returns the index of the first instance of substr in s, or -1 if substr is not present in s.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
