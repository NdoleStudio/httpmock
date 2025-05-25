package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/events"
	"github.com/NdoleStudio/httpmock/pkg/repositories"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
)

// ProjectService is responsible for managing entities.Project
type ProjectService struct {
	service
	logger                           telemetry.Logger
	tracer                           telemetry.Tracer
	repository                       repositories.ProjectRepository
	eventDispatcher                  *EventDispatcher
	projectEndpointRequestRepository repositories.ProjectEndpointRequestRepository
}

// NewProjectService creates a new ProjectService
func NewProjectService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	eventDispatcher *EventDispatcher,
	projectEndpointRequestRepository repositories.ProjectEndpointRequestRepository,
	repository repositories.ProjectRepository,
) (s *ProjectService) {
	return &ProjectService{
		logger:                           logger.WithCodeNamespace(fmt.Sprintf("%T", s)),
		tracer:                           tracer,
		eventDispatcher:                  eventDispatcher,
		projectEndpointRequestRepository: projectEndpointRequestRepository,
		repository:                       repository,
	}
}

// Traffic an entities.Project for an authenticated user
func (service *ProjectService) Traffic(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*repositories.TimeSeriesData, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	traffic, err := service.projectEndpointRequestRepository.GetProjectTraffic(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("could load project traffic for user with ID [%s] and projectID [%s]", userID, projectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg))
	}

	return traffic, nil
}

// Load an entities.Project for an authenticated user
func (service *ProjectService) Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID) (*entities.Project, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	project, err := service.repository.Load(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("could load project for user with ID [%s] and projectID [%s]", userID, projectID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg))
	}

	return project, nil
}

// Index fetches all entities.Project for an authenticated user
func (service *ProjectService) Index(ctx context.Context, userID entities.UserID) ([]*entities.Project, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	projects, err := service.repository.Fetch(ctx, userID)
	if err != nil {
		msg := fmt.Sprintf("could fetch projects for user with ID [%s]", userID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return projects, nil
}

// ProjectCreateParams are the parameters for creating a new project.
type ProjectCreateParams struct {
	Name        string
	Description string
	Subdomain   string
	Source      string
	UserID      entities.UserID
}

// Create a new entities.Project
func (service *ProjectService) Create(ctx context.Context, params *ProjectCreateParams) (*entities.Project, error) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	project := &entities.Project{
		ID:          uuid.New(),
		UserID:      params.UserID,
		Subdomain:   params.Subdomain,
		Name:        params.Name,
		Description: params.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := service.repository.Store(ctx, project); err != nil {
		msg := fmt.Sprintf("could store project [%s] for user with ID [%s]", params.Name, params.UserID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	event, err := service.createEvent(events.ProjectCreated, params.Source, &events.ProjectUpdatedPayload{
		UserID:             project.UserID,
		ProjectID:          project.ID,
		ProjectName:        project.Name,
		ProjectSubdomain:   project.Subdomain,
		ProjectDescription: project.Description,
		ProjectUpdatedAt:   project.UpdatedAt,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for project [%s]", events.ProjectCreated, project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return project, nil
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for project [%s]", event.Type(), project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return project, nil
	}

	return project, nil
}

// ProjectUpdateParams are the parameters for updating a project.
type ProjectUpdateParams struct {
	UserID      entities.UserID
	ProjectID   uuid.UUID
	Subdomain   string
	Name        string
	Description string
	Source      string
}

// Update an entities.Project
func (service *ProjectService) Update(ctx context.Context, params *ProjectUpdateParams) (*entities.Project, error) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	project, err := service.repository.Load(ctx, params.UserID, params.ProjectID)
	if err != nil {
		msg := fmt.Sprintf("cannot load project for user ID [%s] and project [%s]", params.UserID, params.ProjectID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	project.Name = params.Name
	project.UpdatedAt = time.Now().UTC()
	project.Subdomain = params.Subdomain
	project.Description = params.Description

	if err = service.repository.Update(ctx, project); err != nil {
		msg := fmt.Sprintf("could update project [%s] for user with ID [%s]", project.ID, project.UserID)
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	event, err := service.createEvent(events.ProjectUpdated, params.Source, &events.ProjectUpdatedPayload{
		UserID:             project.UserID,
		ProjectID:          project.ID,
		ProjectName:        project.Name,
		ProjectSubdomain:   project.Subdomain,
		ProjectDescription: project.Description,
		ProjectUpdatedAt:   project.UpdatedAt,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for project [%s]", events.ProjectUpdated, project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return project, nil
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for project [%s]", event.Type(), project.ID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return project, nil
	}

	return project, nil
}

// Delete an entities.Project
func (service *ProjectService) Delete(ctx context.Context, source string, userID entities.UserID, projectID uuid.UUID) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if _, err := service.repository.Load(ctx, userID, projectID); err != nil {
		msg := fmt.Sprintf("cannot load project [%s] for user ID [%s]", projectID, userID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	if err := service.repository.Delete(ctx, userID, projectID); err != nil {
		msg := fmt.Sprintf("cannot delete project [%s] for user ID [%s]", projectID, userID)
		return stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	event, err := service.createEvent(events.ProjectDeleted, source, &events.ProjectDeletedPayload{
		UserID:           userID,
		ProjectDeletedAt: time.Now().UTC(),
		ProjectID:        projectID,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot create [%s] event for project [%s]", events.ProjectDeleted, projectID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return nil
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		msg := fmt.Sprintf("cannot dispatch [%s] event for project [%s]", event.Type(), projectID)
		ctxLogger.Error(stacktrace.Propagate(err, msg))
		return nil
	}

	return nil
}
