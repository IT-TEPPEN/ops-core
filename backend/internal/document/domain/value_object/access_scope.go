package value_object

import "errors"

// AccessScope represents the access scope of a document.
type AccessScope string

const (
	// AccessScopePublic represents public access (all users).
	AccessScopePublic AccessScope = "public"
	// AccessScopePrivate represents private access (only owner).
	AccessScopePrivate AccessScope = "private"
	// AccessScopeGroup represents group access (specific groups).
	AccessScopeGroup AccessScope = "group"
	// AccessScopeUser represents user access (specific users).
	AccessScopeUser AccessScope = "user"
)

// NewAccessScope creates a new AccessScope from a string.
func NewAccessScope(scope string) (AccessScope, error) {
	accessScope := AccessScope(scope)
	if !accessScope.IsValid() {
		return "", errors.New("invalid access scope: must be 'public', 'private', 'group', or 'user'")
	}
	return accessScope, nil
}

// IsValid checks if the AccessScope is valid.
func (a AccessScope) IsValid() bool {
	return a == AccessScopePublic || a == AccessScopePrivate || a == AccessScopeGroup || a == AccessScopeUser
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

// IsGroup returns true if the access scope is group.
func (a AccessScope) IsGroup() bool {
	return a == AccessScopeGroup
}

// IsUser returns true if the access scope is user.
func (a AccessScope) IsUser() bool {
	return a == AccessScopeUser
}
