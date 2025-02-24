package events

import (
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/google/uuid"

	"github.com/NdoleStudio/httpmock/pkg/entities"
)

// ProjectEndpointRequest is raised when a new http request is made to an endpoint
const ProjectEndpointRequest = "project.endpoint.request"

// ProjectEndpointRequestPayload stores the data for the ProjectEndpointRequest event
type ProjectEndpointRequestPayload struct {
	UserID                      entities.UserID `json:"user_id"`
	ProjectID                   uuid.UUID       `json:"project_id"`
	ProjectEndpointID           uuid.UUID       `json:"project_endpoint_id"`
	ProjectEndpointRequestID    ulid.ULID       `json:"project_endpoint_request_id"`
	RequestURL                  string          `json:"request_url"`
	RequestMethod               string          `json:"request_method"`
	RequestBody                 *string         `json:"request_body"`
	RequestHeaders              *string         `json:"request_headers"`
	ResponseCode                uint            `json:"response_code"`
	ResponseBody                *string         `json:"response_body"`
	ResponseHeaders             *string         `json:"response_headers"`
	ResponseDelayInMilliseconds uint            `json:"response_delay_in_milliseconds"`
	RequestIPAddress            string          `json:"request_ip_address"`
	Timestamp                   time.Time       `json:"timestamp"`
}
