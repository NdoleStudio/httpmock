package repositories

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/palantir/stacktrace"
	"gorm.io/gorm"
)

// gormProjectEndpointRepository is responsible for persisting entities.ProjectEndpoint
type gormProjectEndpointRequestRepository struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
	db     *gorm.DB
}

// NewGormProjectEndpointRequestRepository creates the GORM version of the ProjectEndpointRequestRepository
func NewGormProjectEndpointRequestRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	db *gorm.DB,
) ProjectEndpointRequestRepository {
	return &gormProjectEndpointRequestRepository{
		logger: logger.WithCodeNamespace(fmt.Sprintf("%T", &gormProjectEndpointRequestRepository{})),
		tracer: tracer,
		db:     db,
	}
}

func (repository *gormProjectEndpointRequestRepository) GetProjectTraffic(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*TimeSeriesData, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	data := new([]*TimeSeriesData)
	err := repository.db.WithContext(ctx).Model(entities.ProjectEndpointRequest{}).
		Select("date_trunc('day', created_at)::date AS timestamp, COUNT(*) AS count").
		Where("user_id = ?", userID).
		Where("project_id = ?", projectID.String()).
		Where(gorm.Expr("created_at >= current_date - interval '30 days'")).
		Group("timestamp").
		Find(&data).Error
	if err != nil {
		msg := fmt.Sprintf("cannot load project traffic for user with ID [%s] and project ID [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return repository.normalizeTimeSeries(*data), nil
}

func (repository *gormProjectEndpointRequestRepository) GetEndpointTraffic(ctx context.Context, userID entities.UserID, endpointID uuid.UUID) ([]*TimeSeriesData, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	data := new([]*TimeSeriesData)
	err := repository.db.WithContext(ctx).Model(entities.ProjectEndpointRequest{}).
		Select("date_trunc('day', created_at)::date AS timestamp, COUNT(*) AS count").
		Where("user_id = ?", userID).
		Where("project_endpoint_id = ?", endpointID.String()).
		Where(gorm.Expr("created_at >= current_date - interval '30 days'")).
		Group("timestamp").
		Find(&data).Error
	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint traffic for user with ID [%s] and project endpoint ID [%s]", userID, endpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return repository.normalizeTimeSeries(*data), nil
}

func (repository *gormProjectEndpointRequestRepository) Index(ctx context.Context, userID entities.UserID, endpointID uuid.UUID, limit uint, previousID *ulid.ULID, nextID *ulid.ULID) ([]*entities.ProjectEndpointRequest, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var endpoints []*entities.ProjectEndpointRequest

	query := repository.db.WithContext(ctx).Where("user_id = ?", userID).Where("project_endpoint_id = ?", endpointID)
	if previousID != nil {
		query = query.Where("id < ?", previousID.String()).Order("id DESC")
	} else if nextID != nil {
		query = query.Where("id > ?", nextID.String()).Order("id ASC")
	} else {
		query = query.Order("id DESC")
	}

	if err := query.Limit(int(limit)).Find(&endpoints).Error; err != nil {
		msg := fmt.Sprintf("cannot load project endpoint requests for user with ID [%s] and endpoint ID [%s]", userID, endpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoints, nil
}

func (repository *gormProjectEndpointRequestRepository) Delete(ctx context.Context, request *entities.ProjectEndpointRequest) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Delete(request).Error; err != nil {
		msg := fmt.Sprintf("cannot delete [%T] with ID [%s] for user [%s]", request, request.ID, request.UserID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectEndpointRequestRepository) Load(ctx context.Context, userID entities.UserID, requestID ulid.ULID) (*entities.ProjectEndpointRequest, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	request := new(entities.ProjectEndpointRequest)
	err := repository.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("id = ?", requestID.String()).
		First(request).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		msg := fmt.Sprintf("request with ID [%s] for userID [%s] does not exist", requestID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}
	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint request for user with ID [%s] and project ID [%s]", userID, requestID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return request, nil
}

func (repository *gormProjectEndpointRequestRepository) Store(ctx context.Context, request *entities.ProjectEndpointRequest) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Create(request).Error; err != nil {
		msg := fmt.Sprintf("cannot save project endpoint request with ID [%s]", request.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *gormProjectEndpointRequestRepository) generateTimeSeries() map[string]*TimeSeriesData {
	series := make(map[string]*TimeSeriesData)
	for i := 0; i < 30; i++ {
		date := time.Now().UTC().AddDate(0, 0, -i)
		series[date.Format("2006-01-02")] = &TimeSeriesData{
			Timestamp: time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -i),
			Count:     0,
		}
	}
	return series
}

func (repository *gormProjectEndpointRequestRepository) normalizeTimeSeries(input []*TimeSeriesData) []*TimeSeriesData {
	series := repository.generateTimeSeries()
	for _, data := range input {
		date := data.Timestamp.Format("2006-01-02")
		if _, ok := series[date]; ok {
			series[date].Count = data.Count
		}
	}

	var result []*TimeSeriesData
	for _, data := range series {
		result = append(result, data)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result
}
