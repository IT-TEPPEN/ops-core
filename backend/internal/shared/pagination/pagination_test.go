package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPage(t *testing.T) {
	page := NewPage()
	assert.Equal(t, 0, page.Offset)
	assert.Equal(t, DefaultPageSize, page.Limit)
}

func TestNewPageFromParams(t *testing.T) {
	tests := []struct {
		name           string
		offset         int
		limit          int
		cursor         string
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "Valid parameters",
			offset:         10,
			limit:          25,
			cursor:         "",
			expectedOffset: 10,
			expectedLimit:  25,
		},
		{
			name:           "Negative offset",
			offset:         -5,
			limit:          20,
			cursor:         "",
			expectedOffset: 0,
			expectedLimit:  20,
		},
		{
			name:           "Zero limit",
			offset:         0,
			limit:          0,
			cursor:         "",
			expectedOffset: 0,
			expectedLimit:  DefaultPageSize,
		},
		{
			name:           "Limit exceeds max",
			offset:         0,
			limit:          200,
			cursor:         "",
			expectedOffset: 0,
			expectedLimit:  MaxPageSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := NewPageFromParams(tt.offset, tt.limit, tt.cursor)
			assert.Equal(t, tt.expectedOffset, page.Offset)
			assert.Equal(t, tt.expectedLimit, page.Limit)
			assert.Equal(t, tt.cursor, page.Cursor)
		})
	}
}

func TestNewPageFromPageNumber(t *testing.T) {
	tests := []struct {
		name           string
		pageNum        int
		pageSize       int
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "First page",
			pageNum:        1,
			pageSize:       20,
			expectedOffset: 0,
			expectedLimit:  20,
		},
		{
			name:           "Second page",
			pageNum:        2,
			pageSize:       20,
			expectedOffset: 20,
			expectedLimit:  20,
		},
		{
			name:           "Page 5",
			pageNum:        5,
			pageSize:       10,
			expectedOffset: 40,
			expectedLimit:  10,
		},
		{
			name:           "Invalid page number",
			pageNum:        0,
			pageSize:       20,
			expectedOffset: 0,
			expectedLimit:  20,
		},
		{
			name:           "Invalid page size",
			pageNum:        1,
			pageSize:       0,
			expectedOffset: 0,
			expectedLimit:  DefaultPageSize,
		},
		{
			name:           "Page size exceeds max",
			pageNum:        1,
			pageSize:       200,
			expectedOffset: 0,
			expectedLimit:  MaxPageSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := NewPageFromPageNumber(tt.pageNum, tt.pageSize)
			assert.Equal(t, tt.expectedOffset, page.Offset)
			assert.Equal(t, tt.expectedLimit, page.Limit)
		})
	}
}

func TestPage_GetPageNumber(t *testing.T) {
	tests := []struct {
		name     string
		offset   int
		limit    int
		expected int
	}{
		{
			name:     "First page",
			offset:   0,
			limit:    20,
			expected: 1,
		},
		{
			name:     "Second page",
			offset:   20,
			limit:    20,
			expected: 2,
		},
		{
			name:     "Fifth page",
			offset:   40,
			limit:    10,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{Offset: tt.offset, Limit: tt.limit}
			assert.Equal(t, tt.expected, page.GetPageNumber())
		})
	}
}

func TestNewPageResult(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	page := NewPageFromPageNumber(2, 3)
	total := int64(10)

	result := NewPageResult(items, total, page)

	assert.Equal(t, items, result.Items)
	assert.Equal(t, total, result.Total)
	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 3, result.PageSize)
	assert.Equal(t, 4, result.TotalPages) // 10 items / 3 per page = 4 pages
	assert.True(t, result.HasNext)        // Page 2 of 4
	assert.True(t, result.HasPrev)        // Not first page
}

func TestNewPageResult_FirstPage(t *testing.T) {
	items := []string{"item1", "item2"}
	page := NewPageFromPageNumber(1, 2)
	total := int64(10)

	result := NewPageResult(items, total, page)

	assert.False(t, result.HasPrev)
	assert.True(t, result.HasNext)
	assert.Equal(t, 1, result.Page)
}

func TestNewPageResult_LastPage(t *testing.T) {
	items := []string{"item1"}
	page := NewPageFromPageNumber(5, 2)
	total := int64(9)

	result := NewPageResult(items, total, page)

	assert.True(t, result.HasPrev)
	assert.False(t, result.HasNext)
	assert.Equal(t, 5, result.Page)
	assert.Equal(t, 5, result.TotalPages)
}

func TestNewCursorPageResult(t *testing.T) {
	items := []string{"item1", "item2", "item3"}
	nextCursor := "cursor123"
	pageSize := 3

	result := NewCursorPageResult(items, true, nextCursor, pageSize)

	assert.Equal(t, items, result.Items)
	assert.True(t, result.HasNext)
	assert.Equal(t, nextCursor, result.NextCursor)
	assert.Equal(t, pageSize, result.PageSize)
	assert.Equal(t, int64(0), result.Total) // Total not used in cursor pagination
}

func TestEncodeDecode_Cursor(t *testing.T) {
	id := "item123"
	timestamp := int64(1234567890)

	encoded, err := EncodeCursor(id, timestamp)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := DecodeCursor(encoded)
	require.NoError(t, err)
	assert.Equal(t, id, decoded.ID)
	assert.Equal(t, timestamp, decoded.Timestamp)
	assert.Equal(t, "next", decoded.Direction)
}

func TestDecodeCursor_Empty(t *testing.T) {
	decoded, err := DecodeCursor("")
	require.NoError(t, err)
	assert.NotNil(t, decoded)
	assert.Empty(t, decoded.ID)
	assert.Equal(t, int64(0), decoded.Timestamp)
}

func TestDecodeCursor_Invalid(t *testing.T) {
	_, err := DecodeCursor("invalid-cursor")
	assert.Error(t, err)
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		defaultVal int
		expected   int
	}{
		{
			name:       "Valid integer",
			input:      "42",
			defaultVal: 10,
			expected:   42,
		},
		{
			name:       "Empty string",
			input:      "",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "Invalid string",
			input:      "abc",
			defaultVal: 10,
			expected:   10,
		},
		{
			name:       "Negative integer",
			input:      "-5",
			defaultVal: 10,
			expected:   -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseInt(tt.input, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{
			name:     "First page",
			page:     1,
			pageSize: 20,
			expected: 0,
		},
		{
			name:     "Second page",
			page:     2,
			pageSize: 20,
			expected: 20,
		},
		{
			name:     "Invalid page",
			page:     0,
			pageSize: 20,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateOffset(tt.page, tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateTotalPages(t *testing.T) {
	tests := []struct {
		name     string
		total    int64
		pageSize int
		expected int
	}{
		{
			name:     "Exact pages",
			total:    100,
			pageSize: 20,
			expected: 5,
		},
		{
			name:     "Partial last page",
			total:    105,
			pageSize: 20,
			expected: 6,
		},
		{
			name:     "Single item",
			total:    1,
			pageSize: 20,
			expected: 1,
		},
		{
			name:     "No items",
			total:    0,
			pageSize: 20,
			expected: 1,
		},
		{
			name:     "Invalid page size",
			total:    100,
			pageSize: 0,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTotalPages(tt.total, tt.pageSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}
