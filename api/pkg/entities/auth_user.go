package entities

// AuthUser is the user gotten from an auth request
type AuthUser struct {
	ID        UserID  `json:"id"`
	Email     string  `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}
