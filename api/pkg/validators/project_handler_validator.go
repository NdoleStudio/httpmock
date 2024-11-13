package validators

import (
	"context"
	"fmt"
	"net/url"

	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/palantir/stacktrace"

	"github.com/NdoleStudio/httpmock/pkg/requests"

	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// ProjectHandlerValidator validates models used in handlers.ProjectHandler
type ProjectHandlerValidator struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	repository repositories.ProjectRepository
}

// NewProjectHandlerValidator creates a new handlers.ProjectHandler validator
func NewProjectHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	repository repositories.ProjectRepository,
) (v *ProjectHandlerValidator) {
	return &ProjectHandlerValidator{
		logger:     logger.WithService(fmt.Sprintf("%T", v)),
		tracer:     tracer,
		repository: repository,
	}
}

// ValidateUpdate validates the requests.ProjectUpdateRequest
func (validator *ProjectHandlerValidator) ValidateUpdate(ctx context.Context, request *requests.ProjectUpdateRequest) url.Values {
	_, span := validator.tracer.Start(ctx)
	defer span.End()

	v := govalidator.New(govalidator.Options{
		Data: request,
		Rules: govalidator.MapData{
			"name": []string{
				"required",
				"min:1",
				"max:30",
			},
			"description": []string{
				"max:500",
			},
		},
	})
	return v.ValidateStruct()
}

// ValidateCreate validates the requests.ProjectCreateRequest
func (validator *ProjectHandlerValidator) ValidateCreate(ctx context.Context, request *requests.ProjectCreateRequest) url.Values {
	_, span, ctxLogger := validator.tracer.StartWithLogger(ctx, validator.logger)
	defer span.End()

	v := govalidator.New(govalidator.Options{
		Data: request,
		Rules: govalidator.MapData{
			"name": []string{
				"required",
				"min:1",
				"max:30",
			},
			"subdomain": []string{
				"required",
				"alpha_dash",
				"min:6",
				"max:30",
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

	exists, err := validator.repository.ExistsWithSubdomain(ctx, request.Subdomain)
	if err != nil {
		msg := fmt.Sprintf("cannot check if the [%s] subdomain has already been taken.", request.Subdomain)
		ctxLogger.Error(validator.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))

		result.Add("subdomain", fmt.Sprintf("We could not check if the [%s] subdomain has already been taken.", request.Subdomain))
		return result
	}

	if exists {
		result.Add("subdomain", fmt.Sprintf("The subdomain [%s] has already been taken.", request.Subdomain))
		return result
	}

	return result
}
