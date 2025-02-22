package requests

import (
	"strings"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/google/uuid"
)

// ProjectEndpointStoreRequest is the payload to update a project endpoint
type ProjectEndpointStoreRequest struct {
	request
	ProjectID string `json:"projectId" swaggerignore:"true"`

	RequestMethod               string `json:"request_method"`
	RequestPath                 string `json:"request_path"`
	ResponseCode                uint   `json:"response_code"`
	ResponseBody                string `json:"response_body"`
	ResponseHeaders             string `json:"response_headers"`
	ResponseDelayInMilliseconds uint   `json:"response_delay_in_milliseconds"`
	Description                 string `json:"description"`
}

// Sanitize the request by stripping whitespaces
func (request *ProjectEndpointStoreRequest) Sanitize() *ProjectEndpointStoreRequest {
	request.RequestMethod = strings.ToUpper(request.sanitizeString(request.RequestMethod))
	request.RequestPath = "/" + strings.TrimLeft(request.sanitizeString(request.RequestPath), "/")
	request.ResponseBody = request.sanitizeString(request.ResponseBody)
	request.ResponseHeaders = request.sanitizeString(request.ResponseHeaders)
	request.Description = request.sanitizeString(request.Description)

	return request
}

// ToProjectEndpointStorePrams creates services.ProjectEndpointStoreParams from ProjectEndpointStoreRequest
func (request *ProjectEndpointStoreRequest) ToProjectEndpointStorePrams(userID entities.UserID) *services.ProjectEndpointStoreParams {
	return &services.ProjectEndpointStoreParams{
		RequestMethod:               request.RequestMethod,
		RequestPath:                 request.RequestPath,
		ResponseCode:                request.ResponseCode,
		ResponseBody:                &request.ResponseBody,
		ResponseHeaders:             &request.ResponseHeaders,
		ResponseDelayInMilliseconds: request.ResponseDelayInMilliseconds,
		Description:                 &request.Description,
		ProjectID:                   uuid.MustParse(request.ProjectID),
		UserID:                      userID,
	}
}
