package requests

import (
	"github.com/oklog/ulid/v2"
)

// ProjectEndpointRequestIndexRequest is the payload fetching entities.ProjectEndpointRequest
type ProjectEndpointRequestIndexRequest struct {
	request

	Prev  string `json:"prev" query:"prev"`
	Next  string `json:"next" query:"next"`
	Limit uint   `json:"limit" query:"limit"`

	ProjectID         string `json:"projectId" swaggerignore:"true"`
	ProjectEndpointID string `json:"projectEndpointId" swaggerignore:"true"`
}

// Sanitize the request by stripping whitespaces
func (input *ProjectEndpointRequestIndexRequest) Sanitize() *ProjectEndpointRequestIndexRequest {
	if input.Limit == 0 {
		input.Limit = 100
	}
	input.Prev = input.sanitizeString(input.Prev)
	return input
}

// PrevID returns the previous ID as a ULID
func (input *ProjectEndpointRequestIndexRequest) PrevID() *ulid.ULID {
	if input.Prev == "" {
		return nil
	}

	id, err := ulid.Parse(input.Prev)
	if err != nil {
		return nil
	}

	return &id
}

// NextID returns the next ID as a ULID
func (input *ProjectEndpointRequestIndexRequest) NextID() *ulid.ULID {
	if input.Prev == "" {
		return nil
	}

	id, err := ulid.Parse(input.Prev)
	if err != nil {
		return nil
	}

	return &id
}
