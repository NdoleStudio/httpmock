package entities

import (
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/google/uuid"
)

// ProjectEndpointRequest is the model for a project endpoint request
type ProjectEndpointRequest struct {
	ID                          ulid.ULID `json:"id" gorm:"primaryKey;type:uuid;" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	ProjectID                   uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	ProjectEndpointID           uuid.UUID `json:"project_endpoint_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID                      UserID    `json:"user_id" example:"user_2oeyIzOf9xxxxxxxxxxxxxx"`
	RequestMethod               string    `json:"request_method" example:"GET" gorm:"type:varchar(7)"`
	RequestURL                  string    `json:"request_url" example:"https://stripe-mock-api.httpmock.dev/v1/products" gorm:"type:varchar(255)"`
	RequestHeaders              *string   `json:"request_headers" example:"[{\"Authorization\":\"Bearer sk_test_4eC39HqLyjWDarjtT1zdp7dc\"}]"`
	RequestBody                 *string   `json:"request_body" example:"{\"name\": \"Product 1\"}"`
	ResponseCode                uint      `json:"response_code" example:"200"`
	ResponseBody                *string   `json:"response_body" example:"{\"message\": \"Hello World\",\"status\": 200}"`
	ResponseHeaders             *string   `json:"response_headers" example:"[{\"Content-Type\":\"application/json\"}]"`
	ResponseDelayInMilliseconds uint      `json:"response_delay_in_milliseconds" example:"1000"`
	CreatedAt                   time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
}
