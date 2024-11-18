package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/google/uuid"

	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/palantir/stacktrace"

	"github.com/NdoleStudio/httpmock/pkg/requests"

	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// ProjectEndpointHandlerValidator validates models used in handlers.ProjectEndpointHandler
type ProjectEndpointHandlerValidator struct {
	validator
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	repository repositories.ProjectEndpointRepository
}

// NewProjectEndpointHandlerValidator creates a new handlers.ProjectEndpointHandler validator
func NewProjectEndpointHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	repository repositories.ProjectEndpointRepository,
) (v *ProjectEndpointHandlerValidator) {
	return &ProjectEndpointHandlerValidator{
		logger:     logger.WithService(fmt.Sprintf("%T", v)),
		tracer:     tracer,
		repository: repository,
	}
}

// ValidateUpdate validates the requests.ProjectUpdateRequest
func (validator *ProjectEndpointHandlerValidator) ValidateUpdate(ctx context.Context, userID entities.UserID, request *requests.ProjectEndpointUpdateRequest) url.Values {
	ctx, span, ctxLogger := validator.tracer.StartWithLogger(ctx, validator.logger)
	defer span.End()

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
			"request_method": []string{
				"required",
				"in:GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD,ANY",
			},
			"request_path": []string{
				"required",
				requestPath,
				"min:1",
				"max:255",
			},
			"response_code": []string{
				"required",
				"min:100",
				"max:600",
			},
			"response_body": []string{
				"max:1000",
			},
			"response_headers": []string{
				requestHeaders,
				"max:500",
			},
			"delay_in_milliseconds": []string{
				"min:0",
				"max:10000",
			},
			"description": []string{
				"max:500",
			},
		},
	})

	result := v.ValidateStruct()
	if len(result) != 0 {
		return result
	}

	endpoint, err := validator.repository.LoadByRequest(ctx, userID, uuid.MustParse(request.ProjectID), request.RequestMethod, request.RequestPath)
	if err != nil && stacktrace.GetCode(err) != repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot check if the [%s %s] request path has already been taken.", request.RequestMethod, request.RequestPath)
		ctxLogger.Error(validator.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))

		result.Add("request_path", fmt.Sprintf("We could not check if the [%s %s] request path has already been taken.", request.RequestMethod, request.RequestPath))
		return result
	}

	if err == nil && endpoint.ID.String() != request.ProjectEndpointID {
		result.Add("request_path", fmt.Sprintf("The request path [%s %s] already exists on this project.", endpoint.RequestMethod, request.RequestPath))
		return result
	}

	return result
}

// ValidateStore validates the requests.ProjectCreateRequest
func (validator *ProjectEndpointHandlerValidator) ValidateStore(ctx context.Context, userID entities.UserID, request *requests.ProjectEndpointStoreRequest) url.Values {
	ctx, span, ctxLogger := validator.tracer.StartWithLogger(ctx, validator.logger)
	defer span.End()

	v := govalidator.New(govalidator.Options{
		Data: request,
		Rules: govalidator.MapData{
			"projectId": []string{
				"required",
				"uuid",
			},
			"request_method": []string{
				"required",
				"in:GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD,ANY",
			},
			"request_path": []string{
				"required",
				requestPath,
				"min:1",
				"max:255",
			},
			"response_code": []string{
				"required",
				"min:100",
				"max:600",
			},
			"response_body": []string{
				"max:1000",
			},
			"response_headers": []string{
				requestHeaders,
				"max:500",
			},
			"delay_in_milliseconds": []string{
				"min:0",
				"max:10000",
			},
			"description": []string{
				"max:500",
			},
		},
	})

	result := v.ValidateStruct()
	if len(result) != 0 {
		return result
	}

	endpoint, err := validator.repository.LoadByRequest(ctx, userID, uuid.MustParse(request.ProjectID), request.RequestMethod, request.RequestPath)
	if err != nil && stacktrace.GetCode(err) != repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot check if the [%s %s] request path has already been taken.", request.RequestMethod, request.RequestPath)
		ctxLogger.Error(validator.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))

		result.Add("request_path", fmt.Sprintf("We could not check if the [%s %s] request path has already been taken.", request.RequestMethod, request.RequestPath))
		return result
	}

	if err == nil {
		result.Add("request_path", fmt.Sprintf("The request path [%s %s] already exists on this project.", endpoint.RequestMethod, request.RequestPath))
		return result
	}

	return result
}
