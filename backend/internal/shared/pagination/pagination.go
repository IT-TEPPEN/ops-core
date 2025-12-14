package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
)

// Page represents pagination parameters.
type Page struct {
	// Offset is the number of items to skip (offset-based pagination)
	Offset int

	// Limit is the maximum number of items to return
	Limit int

	// Cursor is an opaque token for cursor-based pagination
	Cursor string
}

// PageResult represents a paginated result set.
type PageResult struct {
	// Items is the list of items in this page
	Items interface{}

	// Total is the total number of items (for offset-based pagination)
	Total int64

	// HasNext indicates if there is a next page
	HasNext bool

	// HasPrev indicates if there is a previous page
	HasPrev bool

	// NextCursor is the cursor for the next page (cursor-based pagination)
	NextCursor string

	// PrevCursor is the cursor for the previous page (cursor-based pagination)
	PrevCursor string

	// Page is the current page number (1-indexed, for offset-based pagination)
	Page int

	// PageSize is the size of this page
	PageSize int

	// TotalPages is the total number of pages (for offset-based pagination)
	TotalPages int
}

// DefaultPageSize is the default number of items per page.
const DefaultPageSize = 20

// MaxPageSize is the maximum allowed page size.
const MaxPageSize = 100

// NewPage creates a new Page with default values.
func NewPage() *Page {
	return &Page{
		Offset: 0,
		Limit:  DefaultPageSize,
	}
}

// NewPageFromParams creates a Page from query parameters.
func NewPageFromParams(offset, limit int, cursor string) *Page {
	if limit <= 0 {
		limit = DefaultPageSize
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	if offset < 0 {
		offset = 0
	}

	return &Page{
		Offset: offset,
		Limit:  limit,
		Cursor: cursor,
	}
}

// NewPageFromPageNumber creates a Page from a page number (1-indexed).
func NewPageFromPageNumber(page, pageSize int) *Page {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	offset := (page - 1) * pageSize
	return &Page{
		Offset: offset,
		Limit:  pageSize,
	}
}

// GetPageNumber returns the current page number (1-indexed).
func (p *Page) GetPageNumber() int {
	if p.Limit == 0 {
		return 1
	}
	return (p.Offset / p.Limit) + 1
}

// NewPageResult creates a new PageResult for offset-based pagination.
func NewPageResult(items interface{}, total int64, page *Page) *PageResult {
	pageNum := page.GetPageNumber()
	totalPages := int((total + int64(page.Limit) - 1) / int64(page.Limit))
	if totalPages < 1 {
		totalPages = 1
	}

	// Use int64 arithmetic to prevent overflow
	hasNext := int64(page.Offset)+int64(page.Limit) < total
	hasPrev := page.Offset > 0

	return &PageResult{
		Items:      items,
		Total:      total,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
		Page:       pageNum,
		PageSize:   page.Limit,
		TotalPages: totalPages,
	}
}

// NewCursorPageResult creates a new PageResult for cursor-based pagination.
func NewCursorPageResult(items interface{}, hasNext bool, nextCursor string, pageSize int) *PageResult {
	return &PageResult{
		Items:      items,
		HasNext:    hasNext,
		NextCursor: nextCursor,
		PageSize:   pageSize,
	}
}

// Cursor represents a cursor for cursor-based pagination.
type Cursor struct {
	// ID is the last item ID from the previous page
	ID string

	// Timestamp is the last item timestamp from the previous page
	Timestamp int64

	// Direction indicates the pagination direction ("next" or "prev")
	Direction string
}

// EncodeCursor encodes a cursor to a base64 string.
func EncodeCursor(id string, timestamp int64) (string, error) {
	cursor := Cursor{
		ID:        id,
		Timestamp: timestamp,
		Direction: "next",
	}

	data, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor: %w", err)
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

// DecodeCursor decodes a base64 cursor string.
func DecodeCursor(encoded string) (*Cursor, error) {
	if encoded == "" {
		return &Cursor{}, nil
	}

	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cursor: %w", err)
	}

	var cursor Cursor
	if err := json.Unmarshal(data, &cursor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cursor: %w", err)
	}

	return &cursor, nil
}

// ParseInt safely parses an integer from a string.
func ParseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

// CalculateOffset calculates the offset from page number and page size.
func CalculateOffset(page, pageSize int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * pageSize
}

// CalculateTotalPages calculates the total number of pages.
func CalculateTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 1
	}
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if pages < 1 {
		return 1
	}
	return pages
}
