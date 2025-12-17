package db

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIndexDefinitions(t *testing.T) {
	indexes := GetIndexDefinitions()

	// Verify we have the expected number of indexes
	assert.Equal(t, 32, len(indexes), "Should have 32 index definitions")

	// Verify all indexes have required fields
	for _, idx := range indexes {
		assert.NotEmpty(t, idx.Name, "Index name should not be empty")
		assert.NotEmpty(t, idx.Table, "Index table should not be empty")
		assert.NotEmpty(t, idx.Columns, "Index columns should not be empty")
	}

	// Verify unique index names
	nameMap := make(map[string]bool)
	for _, idx := range indexes {
		assert.False(t, nameMap[idx.Name], "Index name %s should be unique", idx.Name)
		nameMap[idx.Name] = true
	}
}

func TestGenerateCreateIndexSQL(t *testing.T) {
	tests := []struct {
		name     string
		index    IndexDefinition
		expected string
	}{
		{
			name: "Simple index",
			index: IndexDefinition{
				Name:    "idx_test",
				Table:   "test_table",
				Columns: []string{"column1"},
				Unique:  false,
			},
			expected: "CREATE INDEX IF NOT EXISTS idx_test ON test_table (column1);",
		},
		{
			name: "Unique index",
			index: IndexDefinition{
				Name:    "idx_test_unique",
				Table:   "test_table",
				Columns: []string{"column1"},
				Unique:  true,
			},
			expected: "CREATE UNIQUE INDEX IF NOT EXISTS idx_test_unique ON test_table (column1);",
		},
		{
			name: "Composite index",
			index: IndexDefinition{
				Name:    "idx_test_composite",
				Table:   "test_table",
				Columns: []string{"column1", "column2", "column3"},
				Unique:  false,
			},
			expected: "CREATE INDEX IF NOT EXISTS idx_test_composite ON test_table (column1, column2, column3);",
		},
		{
			name: "Partial index",
			index: IndexDefinition{
				Name:    "idx_test_partial",
				Table:   "test_table",
				Columns: []string{"column1"},
				Unique:  false,
				Where:   "column1 IS NOT NULL",
			},
			expected: "CREATE INDEX IF NOT EXISTS idx_test_partial ON test_table (column1) WHERE column1 IS NOT NULL;",
		},
		{
			name: "Unique partial index",
			index: IndexDefinition{
				Name:    "idx_test_unique_partial",
				Table:   "test_table",
				Columns: []string{"column1", "column2"},
				Unique:  true,
				Where:   "is_active = true",
			},
			expected: "CREATE UNIQUE INDEX IF NOT EXISTS idx_test_unique_partial ON test_table (column1, column2) WHERE is_active = true;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateCreateIndexSQL(tt.index)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateDropIndexSQL(t *testing.T) {
	tests := []struct {
		name      string
		indexName string
		expected  string
	}{
		{
			name:      "Simple drop",
			indexName: "idx_test",
			expected:  "DROP INDEX IF EXISTS idx_test;",
		},
		{
			name:      "Drop complex name",
			indexName: "idx_documents_repository_id",
			expected:  "DROP INDEX IF EXISTS idx_documents_repository_id;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateDropIndexSQL(tt.indexName)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateCreateIndexSQL_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		index IndexDefinition
	}{
		{
			name: "Invalid index name",
			index: IndexDefinition{
				Name:    "idx-test; DROP TABLE users;--",
				Table:   "test_table",
				Columns: []string{"column1"},
			},
		},
		{
			name: "Invalid table name",
			index: IndexDefinition{
				Name:    "idx_test",
				Table:   "test_table'; DROP TABLE users;--",
				Columns: []string{"column1"},
			},
		},
		{
			name: "Invalid column name",
			index: IndexDefinition{
				Name:    "idx_test",
				Table:   "test_table",
				Columns: []string{"column1; DROP TABLE users;--"},
			},
		},
		{
			name: "Empty columns",
			index: IndexDefinition{
				Name:    "idx_test",
				Table:   "test_table",
				Columns: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateCreateIndexSQL(tt.index)
			assert.Error(t, err, "Should reject invalid input")
		})
	}
}

func TestGenerateDropIndexSQL_InvalidInputs(t *testing.T) {
	tests := []struct {
		name      string
		indexName string
	}{
		{
			name:      "SQL injection attempt",
			indexName: "idx_test; DROP TABLE users;--",
		},
		{
			name:      "Empty name",
			indexName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateDropIndexSQL(tt.indexName)
			assert.Error(t, err, "Should reject invalid input")
		})
	}
}

func TestIndexRationale(t *testing.T) {
	indexes := GetIndexDefinitions()

	// Verify all indexes have a rationale
	for _, idx := range indexes {
		rationale, ok := IndexRationale[idx.Name]
		assert.True(t, ok, "Index %s should have a rationale", idx.Name)
		assert.NotEmpty(t, rationale, "Rationale for %s should not be empty", idx.Name)
	}

	// Verify no extra rationales
	assert.Equal(t, len(indexes), len(IndexRationale), "Number of rationales should match number of indexes")
}

func TestIndexDefinitions_TableCoverage(t *testing.T) {
	indexes := GetIndexDefinitions()
	tableMap := make(map[string]int)

	for _, idx := range indexes {
		tableMap[idx.Table]++
	}

	// Verify we have indexes for key tables
	expectedTables := []string{
		"documents",
		"document_versions",
		"execution_records",
		"execution_steps",
		"attachments",
		"view_history",
		"view_statistics",
		"users",
		"groups",
		"group_members",
	}

	for _, table := range expectedTables {
		assert.Greater(t, tableMap[table], 0, "Should have at least one index for table %s", table)
	}
}

func TestIndexDefinitions_ValidSQL(t *testing.T) {
	indexes := GetIndexDefinitions()

	for _, idx := range indexes {
		sql, err := GenerateCreateIndexSQL(idx)
		require.NoError(t, err, "Should generate valid SQL for index %s", idx.Name)

		// Basic SQL validation
		assert.True(t, strings.HasPrefix(sql, "CREATE"), "SQL should start with CREATE")
		assert.True(t, strings.HasSuffix(sql, ";"), "SQL should end with semicolon")
		assert.Contains(t, sql, "IF NOT EXISTS", "SQL should contain IF NOT EXISTS")
		assert.Contains(t, sql, idx.Name, "SQL should contain index name")
		assert.Contains(t, sql, idx.Table, "SQL should contain table name")

		// Verify all columns are in SQL
		for _, col := range idx.Columns {
			assert.Contains(t, sql, col, "SQL should contain column %s", col)
		}

		// If unique, should contain UNIQUE
		if idx.Unique {
			assert.Contains(t, sql, "UNIQUE", "Unique index should contain UNIQUE keyword")
		}

		// If has WHERE clause, should contain it
		if idx.Where != "" {
			assert.Contains(t, sql, "WHERE", "Partial index should contain WHERE keyword")
			assert.Contains(t, sql, idx.Where, "SQL should contain WHERE clause")
		}
	}
}

func TestIndexDefinitions_UniqueConstraints(t *testing.T) {
	indexes := GetIndexDefinitions()
	uniqueIndexes := []string{
		"idx_document_versions_document_version",
		"idx_execution_steps_record_step",
		"idx_view_statistics_document_id",
		"idx_users_email",
		"idx_users_username",
		"idx_group_members_composite",
	}

	for _, uniqueName := range uniqueIndexes {
		found := false
		for _, idx := range indexes {
			if idx.Name == uniqueName {
				assert.True(t, idx.Unique, "Index %s should be unique", uniqueName)
				found = true
				break
			}
		}
		assert.True(t, found, "Expected unique index %s not found", uniqueName)
	}
}

func TestIndexDefinitions_CompositeIndexes(t *testing.T) {
	indexes := GetIndexDefinitions()

	// Verify we have some composite indexes (multiple columns)
	compositeCount := 0
	for _, idx := range indexes {
		if len(idx.Columns) > 1 {
			compositeCount++
		}
	}

	assert.Greater(t, compositeCount, 0, "Should have at least one composite index")

	// Verify specific composite indexes exist
	compositeIndexes := map[string]int{
		"idx_documents_published":             2, // is_published, access_scope
		"idx_document_versions_document_version": 2, // document_id, version_number
		"idx_execution_records_search":        4, // document_id, executor_id, status, started_at
		"idx_execution_steps_record_step":     2, // execution_record_id, step_number
		"idx_view_history_aggregation":        2, // document_id, viewed_at
		"idx_view_statistics_popularity":      2, // total_views, last_viewed_at
		"idx_group_members_composite":         2, // group_id, user_id
	}

	for indexName, expectedColumns := range compositeIndexes {
		found := false
		for _, idx := range indexes {
			if idx.Name == indexName {
				assert.Equal(t, expectedColumns, len(idx.Columns), 
					"Index %s should have %d columns", indexName, expectedColumns)
				found = true
				break
			}
		}
		assert.True(t, found, "Expected composite index %s not found", indexName)
	}
}

func TestIndexDefinitions_PartialIndexes(t *testing.T) {
	indexes := GetIndexDefinitions()

	// Verify we have at least one partial index
	partialCount := 0
	for _, idx := range indexes {
		if idx.Where != "" {
			partialCount++
		}
	}

	assert.Greater(t, partialCount, 0, "Should have at least one partial index")

	// Verify specific partial index
	for _, idx := range indexes {
		if idx.Name == "idx_documents_published" {
			assert.NotEmpty(t, idx.Where, "idx_documents_published should be a partial index")
			assert.Contains(t, idx.Where, "is_published", "Partial index WHERE clause should reference is_published")
			break
		}
	}
}
