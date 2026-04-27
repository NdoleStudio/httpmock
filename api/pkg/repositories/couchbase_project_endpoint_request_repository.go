package repositories

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/palantir/stacktrace"
)

// couchbaseProjectEndpointRequestRepository is responsible for persisting entities.ProjectEndpointRequest
type couchbaseProjectEndpointRequestRepository struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

// NewCouchbaseProjectEndpointRequestRepository creates the Couchbase version of the ProjectEndpointRequestRepository
func NewCouchbaseProjectEndpointRequestRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	collection *gocb.Collection,
	cluster *gocb.Cluster,
) ProjectEndpointRequestRepository {
	return &couchbaseProjectEndpointRequestRepository{
		logger:     logger.WithCodeNamespace(fmt.Sprintf("%T", &couchbaseProjectEndpointRequestRepository{})),
		tracer:     tracer,
		collection: collection,
		cluster:    cluster,
	}
}

func (repository *couchbaseProjectEndpointRequestRepository) Store(ctx context.Context, request *entities.ProjectEndpointRequest) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Insert(request.ID, request, &gocb.InsertOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot save project endpoint request with ID [%s]", request.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectEndpointRequestRepository) Delete(ctx context.Context, request *entities.ProjectEndpointRequest) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Remove(request.ID, &gocb.RemoveOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot delete [%T] with ID [%s] for user [%s]", request, request.ID, request.UserID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectEndpointRequestRepository) Load(ctx context.Context, userID entities.UserID, requestID ulid.ULID) (*entities.ProjectEndpointRequest, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	result, err := repository.collection.Get(requestID.String(), &gocb.GetOptions{Context: ctx})
	if errors.Is(err, gocb.ErrDocumentNotFound) {
		msg := fmt.Sprintf("request with ID [%s] for userID [%s] does not exist", requestID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}
	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint request for user with ID [%s] and request ID [%s]", userID, requestID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	request := new(entities.ProjectEndpointRequest)
	if err = result.Content(request); err != nil {
		msg := fmt.Sprintf("cannot decode project endpoint request with ID [%s]", requestID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if request.UserID != userID {
		msg := fmt.Sprintf("request with ID [%s] for userID [%s] does not exist", requestID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	return request, nil
}

func (repository *couchbaseProjectEndpointRequestRepository) GetProjectTraffic(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*TimeSeriesData, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	thirtyDaysAgo := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	query := fmt.Sprintf(
		"SELECT DATE_TRUNC_STR(d.created_at, 'day') AS `timestamp`, COUNT(*) AS `count` FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_id = $projectID AND d.created_at >= $thirtyDaysAgo GROUP BY DATE_TRUNC_STR(d.created_at, 'day')",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
		NamedParameters: map[string]interface{}{
			"userID":       string(userID),
			"projectID":    projectID.String(),
			"thirtyDaysAgo": thirtyDaysAgo,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load project traffic for user with ID [%s] and project ID [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	var data []*TimeSeriesData
	for rows.Next() {
		point := new(TimeSeriesData)
		if err = rows.Row(point); err != nil {
			msg := fmt.Sprintf("cannot decode traffic data for user with ID [%s] and project ID [%s]", userID, projectID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		data = append(data, point)
	}

	return repository.normalizeTimeSeries(data), nil
}

func (repository *couchbaseProjectEndpointRequestRepository) GetEndpointTraffic(ctx context.Context, userID entities.UserID, endpointID uuid.UUID) ([]*TimeSeriesData, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	thirtyDaysAgo := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	query := fmt.Sprintf(
		"SELECT DATE_TRUNC_STR(d.created_at, 'day') AS `timestamp`, COUNT(*) AS `count` FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_endpoint_id = $endpointID AND d.created_at >= $thirtyDaysAgo GROUP BY DATE_TRUNC_STR(d.created_at, 'day')",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
		NamedParameters: map[string]interface{}{
			"userID":        string(userID),
			"endpointID":    endpointID.String(),
			"thirtyDaysAgo": thirtyDaysAgo,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint traffic for user with ID [%s] and project endpoint ID [%s]", userID, endpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	var data []*TimeSeriesData
	for rows.Next() {
		point := new(TimeSeriesData)
		if err = rows.Row(point); err != nil {
			msg := fmt.Sprintf("cannot decode traffic data for user with ID [%s] and project endpoint ID [%s]", userID, endpointID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		data = append(data, point)
	}

	return repository.normalizeTimeSeries(data), nil
}

func (repository *couchbaseProjectEndpointRequestRepository) Index(ctx context.Context, userID entities.UserID, endpointID uuid.UUID, limit uint, previousID *ulid.ULID, nextID *ulid.ULID) ([]*entities.ProjectEndpointRequest, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	baseQuery := fmt.Sprintf(
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_endpoint_id = $endpointID",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	params := map[string]interface{}{
		"userID":     string(userID),
		"endpointID": endpointID.String(),
		"limit":      int(limit),
	}

	var query string
	if previousID != nil {
		query = baseQuery + " AND META(d).id < $cursorID ORDER BY META(d).id DESC LIMIT $limit"
		params["cursorID"] = previousID.String()
	} else if nextID != nil {
		query = baseQuery + " AND META(d).id > $cursorID ORDER BY META(d).id ASC LIMIT $limit"
		params["cursorID"] = nextID.String()
	} else {
		query = baseQuery + " ORDER BY META(d).id DESC LIMIT $limit"
	}

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: params,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint requests for user with ID [%s] and endpoint ID [%s]", userID, endpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	var requests []*entities.ProjectEndpointRequest
	for rows.Next() {
		request := new(entities.ProjectEndpointRequest)
		if err = rows.Row(request); err != nil {
			msg := fmt.Sprintf("cannot decode project endpoint request for user with ID [%s] and endpoint ID [%s]", userID, endpointID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (repository *couchbaseProjectEndpointRequestRepository) generateTimeSeries() map[string]*TimeSeriesData {
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

func (repository *couchbaseProjectEndpointRequestRepository) normalizeTimeSeries(input []*TimeSeriesData) []*TimeSeriesData {
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
