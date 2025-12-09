package entity

import (
	"errors"
	"time"

	documentVO "opscore/backend/internal/document/domain/value_object"
	userVO "opscore/backend/internal/user/domain/value_object"
	"opscore/backend/internal/view_history/domain/value_object"
)

// viewHistory represents a document view history record.
type viewHistory struct {
	id         value_object.ViewHistoryID
	documentID documentVO.DocumentID
	userID     userVO.UserID
	viewedAt   time.Time
	viewDuration int // Duration in seconds (0 if not tracked)
}

// ViewHistory is the interface for a view history record.
type ViewHistory interface {
	ID() value_object.ViewHistoryID
	DocumentID() documentVO.DocumentID
	UserID() userVO.UserID
	ViewedAt() time.Time
	ViewDuration() int
}

// NewViewHistory creates a new ViewHistory instance.
func NewViewHistory(
	id value_object.ViewHistoryID,
	documentID documentVO.DocumentID,
	userID userVO.UserID,
	viewedAt time.Time,
) (ViewHistory, error) {
	if id.IsEmpty() {
		return nil, errors.New("view history ID cannot be empty")
	}
	if documentID.IsEmpty() {
		return nil, errors.New("document ID cannot be empty")
	}
	if userID.String() == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	if viewedAt.IsZero() {
		return nil, errors.New("viewed at time cannot be zero")
	}

	return &viewHistory{
		id:         id,
		documentID: documentID,
		userID:     userID,
		viewedAt:   viewedAt,
		viewDuration: 0,
	}, nil
}

// RecordViewHistory creates a new view history record with current time.
func RecordViewHistory(
	documentID documentVO.DocumentID,
	userID userVO.UserID,
) (ViewHistory, error) {
	id := value_object.GenerateViewHistoryID()
	return NewViewHistory(id, documentID, userID, time.Now())
}

// ID returns the view history ID.
func (v *viewHistory) ID() value_object.ViewHistoryID {
	return v.id
}

// DocumentID returns the document ID.
func (v *viewHistory) DocumentID() documentVO.DocumentID {
	return v.documentID
}

// UserID returns the user ID.
func (v *viewHistory) UserID() userVO.UserID {
	return v.userID
}

// ViewedAt returns the time when the document was viewed.
func (v *viewHistory) ViewedAt() time.Time {
	return v.viewedAt
}

// ViewDuration returns the view duration in seconds.
func (v *viewHistory) ViewDuration() int {
	return v.viewDuration
}
