package entities

import (
	"time"
)

// SubscriptionName is the name of the subscription
type SubscriptionName string

// SubscriptionNameFree represents a free subscription
const SubscriptionNameFree = SubscriptionName("free")

// SubscriptionName10kMonthly represents a 10k pro subscription
const SubscriptionName10kMonthly = SubscriptionName("10k-monthly")

// SubscriptionName10kYearly represents a yearly pro subscription
const SubscriptionName10kYearly = SubscriptionName("100k-yearly")

// User stores information about a user
type User struct {
	ID                   UserID           `json:"id" gorm:"primaryKey;type:string;" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	Email                string           `json:"email" example:"name@email.com"`
	FirstName            *string          `json:"first_name" example:"John"`
	LastName             *string          `json:"last_name" example:"Doe"`
	SubscriptionName     SubscriptionName `json:"subscription_name" example:"free"`
	SubscriptionID       string           `json:"subscription_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	SubscriptionStatus   string           `json:"subscription_status" example:"free"`
	SubscriptionRenewsAt *time.Time       `json:"subscription_renews_at" example:"2022-06-05T14:26:02.302718+03:00"`
	SubscriptionEndsAt   *time.Time       `json:"subscription_ends_at" example:"2022-06-05T14:26:02.302718+03:00"`
	CreatedAt            time.Time        `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt            time.Time        `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
