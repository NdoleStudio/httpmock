package entities

import (
	"time"

	"github.com/google/uuid"
)

// Project is a  project belonging to a user
type Project struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID      UserID    `json:"user_id" example:"user_2oeyIzOf9xxxxxxxxxxxxxx"`
	Subdomain   string    `json:"subdomain" example:"stripe-mock-api" gorm:"uniqueIndex;type:varchar(32)"`
	Name        string    `json:"name" example:"Mock Stripe API"`
	Description string    `json:"description" example:"Mock API for an online store for selling shoes"`
	CreatedAt   time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
