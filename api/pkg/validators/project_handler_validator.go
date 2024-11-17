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
	validator
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
	ctx, span, ctxLogger := validator.tracer.StartWithLogger(ctx, validator.logger)
	defer span.End()

	v := govalidator.New(govalidator.Options{
		Data: request,
		Rules: govalidator.MapData{
			"projectId": []string{
				"required",
				"uuid",
			},
			"name": []string{
				"required",
				"min:1",
				"max:30",
			},
			"subdomain": []string{
				"required",
				"alpha_dash",
				"min:7",
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

	project, err := validator.repository.LoadWithSubdomain(ctx, request.Subdomain)
	if err != nil && stacktrace.GetCode(err) != repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot check if the [%s] subdomain has already been taken.", request.Subdomain)
		ctxLogger.Error(validator.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))

		result.Add("subdomain", fmt.Sprintf("We could not check if the [%s] subdomain has already been taken.", request.Subdomain))
		return result
	}

	if err == nil && project.ID.String() != request.ProjectID {
		result.Add("subdomain", fmt.Sprintf("The subdomain [%s] has already been taken.", request.Subdomain))
		return result
	}

	return result
}

// ValidateCreate validates the requests.ProjectCreateRequest
func (validator *ProjectHandlerValidator) ValidateCreate(ctx context.Context, request *requests.ProjectCreateRequest) url.Values {
	ctx, span, ctxLogger := validator.tracer.StartWithLogger(ctx, validator.logger)
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
				"min:7",
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

	_, err := validator.repository.LoadWithSubdomain(ctx, request.Subdomain)
	if err != nil && stacktrace.GetCode(err) != repositories.ErrCodeNotFound {
		msg := fmt.Sprintf("cannot check if the [%s] subdomain has already been taken.", request.Subdomain)
		ctxLogger.Error(validator.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg)))

		result.Add("subdomain", fmt.Sprintf("We could not check if the [%s] subdomain has already been taken.", request.Subdomain))
		return result
	}

	if err != nil {
		result.Add("subdomain", fmt.Sprintf("The subdomain [%s] has already been taken.", request.Subdomain))
		return result
	}

	return result
}
