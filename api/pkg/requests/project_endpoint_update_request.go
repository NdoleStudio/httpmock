package requests

import (
	"strings"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/services"
	"github.com/google/uuid"
)

// ProjectEndpointUpdateRequest is the payload to update a project endpoint
type ProjectEndpointUpdateRequest struct {
	request
	ProjectID         string `json:"projectId" swaggerignore:"true"`
	ProjectEndpointID string `json:"projectEndpointId" swaggerignore:"true"`

	RequestMethod       string `json:"request_method"`
	RequestPath         string `json:"request_path"`
	ResponseCode        uint   `json:"response_code"`
	ResponseBody        string `json:"response_body"`
	ResponseHeaders     string `json:"response_headers"`
	DelayInMilliseconds uint   `json:"delay_in_milliseconds"`
	Description         string `json:"description"`
}

// Sanitize the request by stripping whitespaces
func (request *ProjectEndpointUpdateRequest) Sanitize() *ProjectEndpointUpdateRequest {
	request.RequestMethod = strings.ToUpper(request.sanitizeString(request.RequestMethod))
	request.RequestPath = "/" + strings.TrimLeft(request.sanitizeString(request.RequestPath), "/")
	request.ResponseBody = request.sanitizeString(request.ResponseBody)
	request.ResponseHeaders = request.sanitizeString(request.ResponseHeaders)
	request.Description = request.sanitizeString(request.Description)

	return request
}

// ToProjectEndpointUpdatePrams creates services.ProjectEndpointUpdateParams from ProjectEndpointUpdateRequest
func (request *ProjectEndpointUpdateRequest) ToProjectEndpointUpdatePrams(userID entities.UserID) *services.ProjectEndpointUpdateParams {
	return &services.ProjectEndpointUpdateParams{
		RequestMethod:       request.RequestMethod,
		RequestPath:         request.RequestPath,
		ResponseCode:        request.ResponseCode,
		ResponseBody:        &request.ResponseBody,
		ResponseHeaders:     &request.ResponseHeaders,
		DelayInMilliseconds: request.DelayInMilliseconds,
		Description:         &request.Description,
		ProjectEndpointID:   uuid.MustParse(request.ProjectEndpointID),
		ProjectID:           uuid.MustParse(request.ProjectID),
		UserID:              userID,
	}
}
