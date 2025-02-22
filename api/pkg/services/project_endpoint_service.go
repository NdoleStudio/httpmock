package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
)

// ProjectEndpointService is responsible for managing entities.ProjectEndpoint
type ProjectEndpointService struct {
	service
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	repository repositories.ProjectEndpointRepository
}

// NewProjectEndpointService creates a new ProjectEndpointService
func NewProjectEndpointService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	repository repositories.ProjectEndpointRepository,
) (s *ProjectEndpointService) {
	return &ProjectEndpointService{
		logger:     logger.WithCodeNamespace(fmt.Sprintf("%T", s)),
		tracer:     tracer,
		repository: repository,
	}
}

// Load an entities.Project for an authenticated user
func (service *ProjectEndpointService) Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID, projectEndpointID uuid.UUID) (*entities.ProjectEndpoint, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	project, err := service.repository.Load(ctx, userID, projectID, projectEndpointID)
	if err != nil {
		msg := fmt.Sprintf("could load project endpoint for user with ID [%s], project ID [%s] and endpoint ID [%s]", userID, projectID, projectEndpointID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg))
	}

	return project, nil
}

// Index fetches all entities.Project for an authenticated user
func (service *ProjectEndpointService) Index(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ProjectEndpoint, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	endpoints, err := service.repository.Fetch(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("could load project endpoint for user with ID [%s] abd project ID [%s]", userID, projectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoints, nil
}

// ProjectEndpointStoreParams are the parameters for creating a new entities.ProjectEndpoint.
type ProjectEndpointStoreParams struct {
	RequestMethod       string
	RequestPath         string
	ResponseCode        uint
	ResponseBody        *string
	ResponseHeaders     *string
	DelayInMilliseconds uint
	Description         *string

	ProjectID uuid.UUID
	UserID    entities.UserID
}

// Store a new entities.Project
func (service *ProjectEndpointService) Store(ctx context.Context, project *entities.Project, params *ProjectEndpointStoreParams) (*entities.ProjectEndpoint, error) {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	endpoint := &entities.ProjectEndpoint{
		ID:              uuid.New(),
		UserID:          params.UserID,
		ProjectID:       params.ProjectID,
		RequestMethod:   params.RequestMethod,
		RequestPath:     params.RequestPath,
		ResponseCode:    params.ResponseCode,
		ResponseBody:    params.ResponseBody,
		ResponseHeaders: params.ResponseHeaders,
		Subdomain:       project.Subdomain,
		Description:     params.Description,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := service.repository.Store(ctx, endpoint); err != nil {
		msg := fmt.Sprintf("could store project endpoint for user with ID [%s] and project ID [%s]", params.UserID, params.ProjectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoint, nil
}

// ProjectEndpointUpdateParams are the parameters for updating a project endpoint.
type ProjectEndpointUpdateParams struct {
	RequestMethod       string
	RequestPath         string
	ResponseCode        uint
	ResponseBody        *string
	ResponseHeaders     *string
	DelayInMilliseconds uint
	Description         *string

	ProjectEndpointID uuid.UUID
	ProjectID         uuid.UUID
	UserID            entities.UserID
}

// Update an entities.Project
func (service *ProjectEndpointService) Update(ctx context.Context, params *ProjectEndpointUpdateParams) (*entities.ProjectEndpoint, error) {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	endpoint, err := service.repository.Load(ctx, params.UserID, params.ProjectID, params.ProjectEndpointID)
	if err != nil {
		msg := fmt.Sprintf("cannot load endpoint for user ID [%s], project ID [%s] and endpoint ID [%s]", params.UserID, params.ProjectID, params.ProjectEndpointID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	endpoint.RequestMethod = params.RequestMethod
	endpoint.RequestPath = params.RequestPath
	endpoint.ResponseCode = params.ResponseCode
	endpoint.ResponseBody = params.ResponseBody
	endpoint.ResponseHeaders = params.ResponseHeaders
	endpoint.DelayInMilliseconds = params.DelayInMilliseconds
	endpoint.Description = params.Description

	endpoint.UpdatedAt = time.Now().UTC()

	if err = service.repository.Update(ctx, endpoint); err != nil {
		msg := fmt.Sprintf("cannot update endpoint for user ID [%s], project ID [%s] and endpoint ID [%s]", params.UserID, params.ProjectID, params.ProjectEndpointID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoint, nil
}

// Delete an entities.Project
func (service *ProjectEndpointService) Delete(ctx context.Context, userID entities.UserID, projectID uuid.UUID, projectEndpoint uuid.UUID) error {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	endpoint, err := service.repository.Load(ctx, userID, projectID, projectEndpoint)
	if err != nil {
		msg := fmt.Sprintf("cannot load endpoint with ID [%s] and project ID [%s] for user ID [%s]", projectEndpoint, projectID, userID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	if err = service.repository.Delete(ctx, endpoint); err != nil {
		msg := fmt.Sprintf("cannot delete endpoint with ID [%s] and project ID [%s] for user ID [%s]", projectEndpoint, projectID, userID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	return nil
}
