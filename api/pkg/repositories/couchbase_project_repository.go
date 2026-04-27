package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
)

// couchbaseProjectRepository is responsible for persisting entities.Project
type couchbaseProjectRepository struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

// NewCouchbaseProjectRepository creates the Couchbase version of the ProjectRepository
func NewCouchbaseProjectRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	collection *gocb.Collection,
	cluster *gocb.Cluster,
) ProjectRepository {
	return &couchbaseProjectRepository{
		logger:     logger.WithCodeNamespace(fmt.Sprintf("%T", &couchbaseProjectRepository{})),
		tracer:     tracer,
		collection: collection,
		cluster:    cluster,
	}
}

func (repository *couchbaseProjectRepository) Store(ctx context.Context, project *entities.Project) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Insert(project.ID.String(), project, &gocb.InsertOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot save project with ID [%s]", project.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectRepository) Update(ctx context.Context, project *entities.Project) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Upsert(project.ID.String(), project, &gocb.UpsertOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot update project with ID [%s]", project.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectRepository) Fetch(ctx context.Context, userID entities.UserID) ([]*entities.Project, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := fmt.Sprintf(
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID ORDER BY d.created_at DESC",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: map[string]interface{}{"userID": string(userID)},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load projects for user with ID [%s]", userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	var projects []*entities.Project
	for rows.Next() {
		project := new(entities.Project)
		if err = rows.Row(project); err != nil {
			msg := fmt.Sprintf("cannot decode project for user with ID [%s]", userID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		projects = append(projects, project)
	}

	return projects, nil
}

func (repository *couchbaseProjectRepository) Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID) (*entities.Project, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	result, err := repository.collection.Get(projectID.String(), &gocb.GetOptions{Context: ctx})
	if errors.Is(err, gocb.ErrDocumentNotFound) {
		msg := fmt.Sprintf("project with ID [%s] does not exist", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}
	if err != nil {
		msg := fmt.Sprintf("cannot load project with ID [%s]", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	project := new(entities.Project)
	if err = result.Content(project); err != nil {
		msg := fmt.Sprintf("cannot decode project with ID [%s]", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if project.UserID != userID {
		msg := fmt.Sprintf("project with ID [%s] does not exist", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	return project, nil
}

func (repository *couchbaseProjectRepository) Delete(ctx context.Context, userID entities.UserID, projectID uuid.UUID) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	// Verify ownership before deleting
	_, err := repository.Load(ctx, userID, projectID)
	if err != nil {
		return err
	}

	_, err = repository.collection.Remove(projectID.String(), &gocb.RemoveOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot delete project with ID [%s] for user [%s]", projectID, userID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseProjectRepository) LoadWithSubdomain(ctx context.Context, subdomain string) (*entities.Project, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := fmt.Sprintf(
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.subdomain = $subdomain",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: map[string]interface{}{"subdomain": subdomain},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load a project with subdomain [%s]", subdomain)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	project := new(entities.Project)
	if !rows.Next() {
		msg := fmt.Sprintf("project not found with subdomain [%s]", subdomain)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	if err = rows.Row(project); err != nil {
		msg := fmt.Sprintf("cannot decode project with subdomain [%s]", subdomain)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return project, nil
}
