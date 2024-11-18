package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// gormProjectEndpointRepository is responsible for persisting entities.ProjectEndpoint
type gormProjectEndpointRepository struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
	db     *gorm.DB
}

// NewGormProjectEndpointRepository creates the GORM version of the ProjectEndpointRepository
func NewGormProjectEndpointRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) ProjectEndpointRepository {
	return &gormProjectEndpointRepository{
		logger: logger.WithService(fmt.Sprintf("%T", &gormProjectEndpointRepository{})),
		tracer: tracer,
		db:     db,
	}
}

func (repository *gormProjectEndpointRepository) LoadByRequest(ctx context.Context, userID entities.UserID, projectID uuid.UUID, requestMethod, requestPath string) (*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	endpoint := new(entities.ProjectEndpoint)

	query := repository.db.WithContext(ctx).
		Model(endpoint).
		Where("user_id = ?", userID).
		Where("project_id = ?", projectID.String()).
		Where("request_path = ?", requestPath)

	if requestMethod != "ANY" {
		query.Where(repository.db.Where("request_method = ?", requestMethod).Or("request_method = ?", "ANY"))
	}

	err := query.First(endpoint).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("endpoint not found with project ID [%s], request method [%s] and request path [%s]", projectID, requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("endpoint not found with project ID [%s], request method [%s] and request path [%s]", projectID, requestMethod, requestPath)
		return endpoint, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoint, nil
}

func (repository *gormProjectEndpointRepository) Delete(ctx context.Context, userID entities.UserID, projectID uuid.UUID, projectEndpointID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("project_id = ?", projectID).
		Where("id = ?", projectEndpointID).
		Delete(&entities.Project{}).
		Error
	if err != nil {
		msg := fmt.Sprintf("cannot save project endpoint with ID [%s] for user [%s]", projectEndpointID, userID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	return nil
}

func (repository *gormProjectEndpointRepository) Store(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Create(endpoint).Error; err != nil {
		msg := fmt.Sprintf("cannot save project endpoint with ID [%s]", endpoint.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectEndpointRepository) Update(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Save(endpoint).Error; err != nil {
		msg := fmt.Sprintf("cannot update project endpoint with ID [%s]", endpoint.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectEndpointRepository) Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var endpoints []*entities.ProjectEndpoint
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("project_id = ?", projectID).
		Find(&endpoints).
		Error
	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint for user with ID [%s] and project ID [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoints, nil
}

func (repository *gormProjectEndpointRepository) Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID, projectEndpointID uuid.UUID) (*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	endpoint := new(entities.ProjectEndpoint)
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("project_id = ?", projectID).
		Where("id = ?", projectEndpointID).
		First(endpoint).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("project with ID [%s] does not exist", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("cannot load project wndpoint with ID [%s]", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoint, nil
}
