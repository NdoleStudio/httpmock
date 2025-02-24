package entities

// UserID is the ID of a user
type UserID string

// String returns the string representation of the UserID
func (id UserID) String() string {
	return string(id)
}
