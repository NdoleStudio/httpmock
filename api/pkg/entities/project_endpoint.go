package entities

import (
	"time"

	"github.com/google/uuid"
)

// ProjectEndpoint is an endpoint belonging to a project
type ProjectEndpoint struct {
	ID                  uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	ProjectID           uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID              UserID    `json:"user_id" example:"user_2oeyIzOf9xxxxxxxxxxxxxx"`
	RequestMethod       string    `json:"request_method" example:"GET" gorm:"type:varchar(7)"`
	RequestPath         string    `json:"request_path" example:"/v1/products" gorm:"type:varchar(255)"`
	ResponseCode        uint      `json:"response_code" example:"200"`
	ResponseBody        *string   `json:"response_body" example:"{\"message\": \"Hello World\",\"status\": 200}"`
	ResponseHeaders     *string   `json:"response_headers" example:"[{\"Content-Type\":\"application/json\"}]"`
	DelayInMilliseconds uint      `json:"delay_in_milliseconds" example:"1000"`
	Description         *string   `json:"description" example:"Mock API for an online store for the /v1/products endpoint"`
	CreatedAt           time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt           time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
