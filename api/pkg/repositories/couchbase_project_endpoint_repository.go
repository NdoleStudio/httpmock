package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
)

// couchbaseProjectEndpointRepository is responsible for persisting entities.ProjectEndpoint
type couchbaseProjectEndpointRepository struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

// NewCouchbaseProjectEndpointRepository creates the Couchbase version of the ProjectEndpointRepository
func NewCouchbaseProjectEndpointRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	collection *gocb.Collection,
	cluster *gocb.Cluster,
) ProjectEndpointRepository {
	return &couchbaseProjectEndpointRepository{
		logger:     logger.WithCodeNamespace(fmt.Sprintf("%T", &couchbaseProjectEndpointRepository{})),
		tracer:     tracer,
		collection: collection,
		cluster:    cluster,
	}
}

func (repository *couchbaseProjectEndpointRepository) Store(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Insert(endpoint.ID.String(), endpoint, &gocb.InsertOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot save project endpoint with ID [%s]", endpoint.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectEndpointRepository) Update(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Upsert(endpoint.ID.String(), endpoint, &gocb.UpsertOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot update project endpoint with ID [%s]", endpoint.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectEndpointRepository) IncreaseRequestCount(ctx context.Context, projectEndpointID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.MutateIn(projectEndpointID.String(), []gocb.MutateInSpec{
		gocb.IncrementSpec("request_count", int64(1), &gocb.CounterSpecOptions{}),
	}, &gocb.MutateInOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot increase request_count [%T] with ID [%s]", &entities.ProjectEndpoint{}, projectEndpointID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectEndpointRepository) DecreaseRequestCount(ctx context.Context, projectEndpointID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.MutateIn(projectEndpointID.String(), []gocb.MutateInSpec{
		gocb.DecrementSpec("request_count", int64(1), &gocb.CounterSpecOptions{}),
	}, &gocb.MutateInOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot decrease request_count [%T] with ID [%s]", &entities.ProjectEndpoint{}, projectEndpointID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectEndpointRepository) UpdateSubdomain(ctx context.Context, subdomain string, projectID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := fmt.Sprintf(
		"UPDATE `%s`.`%s`.`%s` SET project_subdomain = $subdomain WHERE project_id = $projectID",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	_, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
		NamedParameters: map[string]interface{}{
			"subdomain": subdomain,
			"projectID": projectID.String(),
		},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot update [project_subdomain] for [%T] with project ID [%s]", &entities.ProjectEndpoint{}, projectID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectEndpointRepository) Fetch(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := fmt.Sprintf(
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_id = $projectID ORDER BY d.created_at DESC",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
		NamedParameters: map[string]interface{}{
			"userID":    string(userID),
			"projectID": projectID.String(),
		},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint for user with ID [%s] and project ID [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	var endpoints []*entities.ProjectEndpoint
	for rows.Next() {
		endpoint := new(entities.ProjectEndpoint)
		if err = rows.Row(endpoint); err != nil {
			msg := fmt.Sprintf("cannot decode project endpoint for user with ID [%s] and project ID [%s]", userID, projectID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}

func (repository *couchbaseProjectEndpointRepository) Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID, projectEndpointID uuid.UUID) (*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	result, err := repository.collection.Get(projectEndpointID.String(), &gocb.GetOptions{Context: ctx})
	if errors.Is(err, gocb.ErrDocumentNotFound) {
		msg := fmt.Sprintf("project endpoint with ID [%s] does not exist", projectEndpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}
	if err != nil {
		msg := fmt.Sprintf("cannot load project endpoint with ID [%s]", projectEndpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	endpoint := new(entities.ProjectEndpoint)
	if err = result.Content(endpoint); err != nil {
		msg := fmt.Sprintf("cannot decode project endpoint with ID [%s]", projectEndpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if endpoint.UserID != userID || endpoint.ProjectID != projectID {
		msg := fmt.Sprintf("project endpoint with ID [%s] does not exist", projectEndpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	return endpoint, nil
}

func (repository *couchbaseProjectEndpointRepository) Delete(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Remove(endpoint.ID.String(), &gocb.RemoveOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot delete [%T] with ID [%s] for user [%s]", endpoint, endpoint.ID, endpoint.UserID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectEndpointRepository) LoadByRequest(ctx context.Context, subdomain, requestMethod, requestPath string) (*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := fmt.Sprintf(
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.project_subdomain = $subdomain AND d.request_path = $path AND (d.request_method = $method OR d.request_method = 'ANY')",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
		NamedParameters: map[string]interface{}{
			"subdomain": subdomain,
			"path":      requestPath,
			"method":    strings.ToUpper(requestMethod),
		},
	})
	if err != nil {
		msg := fmt.Sprintf("endpoint not found with request method [%s] and request path [%s]", requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	endpoint := new(entities.ProjectEndpoint)
	if !rows.Next() {
		msg := fmt.Sprintf("endpoint not found with request method [%s] and request path [%s]", requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	if err = rows.Row(endpoint); err != nil {
		msg := fmt.Sprintf("cannot decode endpoint with request method [%s] and request path [%s]", requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoint, nil
}

func (repository *couchbaseProjectEndpointRepository) LoadByRequestForUser(ctx context.Context, userID entities.UserID, projectID uuid.UUID, requestMethod, requestPath string) (*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var query string
	params := map[string]interface{}{
		"userID":    string(userID),
		"projectID": projectID.String(),
		"path":      requestPath,
	}

	if requestMethod == "ANY" {
		query = fmt.Sprintf(
			"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_id = $projectID AND d.request_path = $path",
			repository.collection.Bucket().Name(),
			repository.collection.ScopeName(),
			repository.collection.Name(),
		)
	} else {
		query = fmt.Sprintf(
			"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_id = $projectID AND d.request_path = $path AND (d.request_method = $method OR d.request_method = 'ANY')",
			repository.collection.Bucket().Name(),
			repository.collection.ScopeName(),
			repository.collection.Name(),
		)
		params["method"] = requestMethod
	}

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: params,
	})
	if err != nil {
		msg := fmt.Sprintf("endpoint not found with project ID [%s], request method [%s] and request path [%s]", projectID, requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	endpoint := new(entities.ProjectEndpoint)
	if !rows.Next() {
		msg := fmt.Sprintf("endpoint not found with project ID [%s], request method [%s] and request path [%s]", projectID, requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	if err = rows.Row(endpoint); err != nil {
		msg := fmt.Sprintf("cannot decode endpoint with project ID [%s], request method [%s] and request path [%s]", projectID, requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return endpoint, nil
}
