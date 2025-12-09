package entity

import (
	"errors"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	"opscore/backend/internal/view_statistics/domain/value_object"
)

// viewStatistics represents aggregated view statistics for a document.
type viewStatistics struct {
	id                 value_object.ViewStatisticsID
	documentID         documentVO.DocumentID
	totalViews         int64
	uniqueViewers      int64
	lastViewedAt       time.Time
	averageViewDuration int // in seconds
	createdAt          time.Time
	updatedAt          time.Time
}

// ViewStatistics is the interface for view statistics (aggregate root).
type ViewStatistics interface {
	ID() value_object.ViewStatisticsID
	DocumentID() documentVO.DocumentID
	TotalViews() int64
	UniqueViewers() int64
	LastViewedAt() time.Time
	AverageViewDuration() int
	CreatedAt() time.Time
	UpdatedAt() time.Time

	// Behaviors
	IncrementView(isUniqueViewer bool) error
	UpdateLastViewedAt(viewedAt time.Time) error
	UpdateAverageViewDuration(duration int) error
}

// NewViewStatistics creates a new ViewStatistics instance.
func NewViewStatistics(
	id value_object.ViewStatisticsID,
	documentID documentVO.DocumentID,
) (ViewStatistics, error) {
	if id.IsEmpty() {
		return nil, errors.New("view statistics ID cannot be empty")
	}
	if documentID.IsEmpty() {
		return nil, errors.New("document ID cannot be empty")
	}

	now := time.Now()
	return &viewStatistics{
		id:                 id,
		documentID:         documentID,
		totalViews:         0,
		uniqueViewers:      0,
		lastViewedAt:       time.Time{},
		averageViewDuration: 0,
		createdAt:          now,
		updatedAt:          now,
	}, nil
}

// ID returns the view statistics ID.
func (v *viewStatistics) ID() value_object.ViewStatisticsID {
	return v.id
}

// DocumentID returns the document ID.
func (v *viewStatistics) DocumentID() documentVO.DocumentID {
	return v.documentID
}

// TotalViews returns the total number of views.
func (v *viewStatistics) TotalViews() int64 {
	return v.totalViews
}

// UniqueViewers returns the number of unique viewers.
func (v *viewStatistics) UniqueViewers() int64 {
	return v.uniqueViewers
}

// LastViewedAt returns the last viewed timestamp.
func (v *viewStatistics) LastViewedAt() time.Time {
	return v.lastViewedAt
}

// AverageViewDuration returns the average view duration in seconds.
func (v *viewStatistics) AverageViewDuration() int {
	return v.averageViewDuration
}

// CreatedAt returns the creation timestamp.
func (v *viewStatistics) CreatedAt() time.Time {
	return v.createdAt
}

// UpdatedAt returns the last update timestamp.
func (v *viewStatistics) UpdatedAt() time.Time {
	return v.updatedAt
}

// IncrementView increments the view count.
func (v *viewStatistics) IncrementView(isUniqueViewer bool) error {
	v.totalViews++
	if isUniqueViewer {
		v.uniqueViewers++
	}
	v.updatedAt = time.Now()
	return nil
}

// UpdateLastViewedAt updates the last viewed timestamp.
func (v *viewStatistics) UpdateLastViewedAt(viewedAt time.Time) error {
	if viewedAt.IsZero() {
		return errors.New("viewed at time cannot be zero")
	}
	v.lastViewedAt = viewedAt
	v.updatedAt = time.Now()
	return nil
}

// UpdateAverageViewDuration updates the average view duration.
func (v *viewStatistics) UpdateAverageViewDuration(duration int) error {
	if duration < 0 {
		return errors.New("duration cannot be negative")
	}
	v.averageViewDuration = duration
	v.updatedAt = time.Now()
	return nil
}
