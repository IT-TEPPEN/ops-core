package value_object

import "errors"

// AccessScope represents the access scope of an execution record.
// Public records are accessible to all users.
// Private records are only accessible to the executor, but can be shared with specific users or groups.
type AccessScope string

const (
	// AccessScopePublic represents public access (all users).
	AccessScopePublic AccessScope = "public"
	// AccessScopePrivate represents private access (executor only, but can be shared).
	AccessScopePrivate AccessScope = "private"
)

// NewAccessScope creates a new AccessScope from a string.
func NewAccessScope(scope string) (AccessScope, error) {
	accessScope := AccessScope(scope)
	if !accessScope.IsValid() {
		return "", errors.New("invalid access scope: must be 'public' or 'private'")
	}
	return accessScope, nil
}

// IsValid checks if the AccessScope is valid.
func (a AccessScope) IsValid() bool {
	return a == AccessScopePublic || a == AccessScopePrivate
}

// String returns the string representation of AccessScope.
func (a AccessScope) String() string {
	return string(a)
}

// Equals checks if two AccessScopes are equal.
func (a AccessScope) Equals(other AccessScope) bool {
	return a == other
}

// IsPublic returns true if the access scope is public.
func (a AccessScope) IsPublic() bool {
	return a == AccessScopePublic
}

// IsPrivate returns true if the access scope is private.
func (a AccessScope) IsPrivate() bool {
	return a == AccessScopePrivate
}
