package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto/v2"

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
	cache  *ristretto.Cache[string, *entities.ProjectEndpoint]
	db     *gorm.DB
}

func (repository *gormProjectEndpointRepository) UpdateSubdomain(ctx context.Context, subdomain string, projectID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	repository.cache.Clear()
	err := repository.db.WithContext(ctx).
		Model(&entities.ProjectEndpoint{}).
		Where("project_id = ?", projectID).
		UpdateColumn("project_subdomain", subdomain).Error
	if err != nil {
		msg := fmt.Sprintf("cannot update [project_subdomain] for [%T] with project ID [%s]", &entities.ProjectEndpoint{}, projectID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

// NewGormProjectEndpointRepository creates the GORM version of the ProjectEndpointRepository
func NewGormProjectEndpointRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) ProjectEndpointRepository {
	return &gormProjectEndpointRepository{
		logger: logger.WithCodeNamespace(fmt.Sprintf("%T", &gormProjectEndpointRepository{})),
		tracer: tracer,
		db:     db,
	}
}

func (repository *gormProjectEndpointRepository) RegisterRequest(ctx context.Context, projectEndpointID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	err := repository.db.WithContext(ctx).
		Model(&entities.ProjectEndpoint{}).
		Where("id = ?", projectEndpointID).
		UpdateColumn("request_count", gorm.Expr("request_count + ?", 1)).Error
	if err != nil {
		msg := fmt.Sprintf("cannot update request_count [%T] with ID [%s]", &entities.ProjectEndpoint{}, projectEndpointID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectEndpointRepository) LoadByRequest(ctx context.Context, subdomain, requestMethod, requestPath string) (*entities.ProjectEndpoint, error) {
	ctx, span, ctxLogger := repository.tracer.StartWithLogger(ctx, repository.logger)
	defer span.End()

	key := repository.cacheKey(subdomain, requestMethod, requestPath)
	if endpoint, ok := repository.cache.Get(key); ok {
		ctxLogger.Info(fmt.Sprintf("[%T] found in cache with ID [%s] for request with key [%s]", endpoint, endpoint.ID, key))
		return endpoint, nil
	}

	endpoint := new(entities.ProjectEndpoint)
	err := repository.db.WithContext(ctx).
		Model(endpoint).
		Where("subdomain = ?", subdomain).
		Where("request_path = ?", requestPath).
		Where(repository.db.Where("request_method = ?", strings.ToUpper(requestMethod)).Or("request_method = ?", "ANY")).
		First(endpoint).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("endpoint not found with request method [%s] and request path [%s]", requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}

	if err != nil {
		msg := fmt.Sprintf("endpoint not found with request method [%s] and request path [%s]", requestMethod, requestPath)
		return endpoint, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if ok := repository.cache.SetWithTTL(key, endpoint, 1, 2*time.Hour); !ok {
		msg := fmt.Sprintf("cannot cache [%T] with ID [%s]", endpoint, endpoint.ID)
		ctxLogger.Error(repository.tracer.WrapErrorSpan(span, stacktrace.NewError(msg)))
	}

	return endpoint, nil
}

func (repository *gormProjectEndpointRepository) LoadByRequestForUser(ctx context.Context, userID entities.UserID, projectID uuid.UUID, requestMethod, requestPath string) (*entities.ProjectEndpoint, error) {
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

func (repository *gormProjectEndpointRepository) Delete(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	repository.cache.Del(repository.cacheKeyFromEndpoint(endpoint))
	if err := repository.db.WithContext(ctx).Delete(endpoint).Error; err != nil {
		msg := fmt.Sprintf("cannot delete [%T] with ID [%s] for user [%s]", endpoint, endpoint.ID, endpoint.UserID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectEndpointRepository) Store(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	repository.cache.Del(repository.cacheKeyFromEndpoint(endpoint))
	if err := repository.db.WithContext(ctx).Create(endpoint).Error; err != nil {
		msg := fmt.Sprintf("cannot save project endpoint with ID [%s]", endpoint.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectEndpointRepository) Update(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	repository.cache.Del(repository.cacheKeyFromEndpoint(endpoint))
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

func (repository *gormProjectEndpointRepository) cacheKeyFromEndpoint(endpoint *entities.ProjectEndpoint) string {
	return repository.cacheKey(endpoint.ProjectSubdomain, endpoint.RequestMethod, endpoint.RequestPath)
}

func (repository *gormProjectEndpointRepository) cacheKey(subdomain, method, path string) string {
	return fmt.Sprintf("[%s] [%s] [%s]", subdomain, method, path)
}
