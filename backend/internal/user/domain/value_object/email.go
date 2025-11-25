package value_object

import (
	"errors"
	"regexp"
	"strings"
)

// Email represents a validated email address
type Email struct {
	value string
}

// Email validation regex pattern
// This is a simplified pattern that covers most common email formats
var emailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email with validation
func NewEmail(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	
	normalized := strings.TrimSpace(strings.ToLower(email))
	
	if !emailPattern.MatchString(normalized) {
		return Email{}, errors.New("invalid email format")
	}
	
	return Email{value: normalized}, nil
}

// String returns the string representation of Email
func (e Email) String() string {
	return e.value
}

// Equals checks if two Emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
