package validators

import (
	"fmt"
	"net/url"

	"github.com/oklog/ulid/v2"

	"github.com/NdoleStudio/httpmock/pkg/requests"

	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// ProjectEndpointRequestHandlerValidator validates models used in handlers.ProjectEndpointRequestHandler
type ProjectEndpointRequestHandlerValidator struct {
	validator
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewProjectEndpointRequestHandlerValidator creates a new handlers.ProjectEndpointRequestHandler validator
func NewProjectEndpointRequestHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (v *ProjectEndpointRequestHandlerValidator) {
	return &ProjectEndpointRequestHandlerValidator{
		logger: logger.WithCodeNamespace(fmt.Sprintf("%T", v)),
		tracer: tracer,
	}
}

// ValidateIndex validates the requests.ProjectEndpointRequestIndexRequest
func (validator *ProjectEndpointRequestHandlerValidator) ValidateIndex(request *requests.ProjectEndpointRequestIndexRequest) url.Values {
	v := govalidator.New(govalidator.Options{
		Data: request,
		Rules: govalidator.MapData{
			"projectId": []string{
				"required",
				"uuid",
			},
			"projectEndpointId": []string{
				"required",
				"uuid",
			},
			"limit": []string{
				"required",
				"min:1",
				"max:100",
			},
		},
	})

	validationErrors := v.ValidateStruct()
	if len(validationErrors) > 0 || request.Prev == "" {
		return validationErrors
	}

	if _, err := ulid.Parse(request.Prev); request.Prev != "" && err != nil {
		validationErrors["prev"] = []string{fmt.Sprintf("The prev query param [%s] must be a valid ULID https://github.com/ulid/spec", request.Prev)}
	}

	if _, err := ulid.Parse(request.Next); request.Next != "" && err != nil {
		validationErrors["next"] = []string{fmt.Sprintf("The next query param [%s] must be a valid ULID https://github.com/ulid/spec", request.Next)}
	}

	return validationErrors
}
