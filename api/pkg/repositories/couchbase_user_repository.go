package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NdoleStudio/httpmock/pkg/entities"
	"github.com/NdoleStudio/httpmock/pkg/telemetry"
	"github.com/couchbase/gocb/v2"
	"github.com/palantir/stacktrace"
)

// couchbaseUserRepository is responsible for persisting entities.User
type couchbaseUserRepository struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

// NewCouchbaseUserRepository creates the Couchbase version of the UserRepository
func NewCouchbaseUserRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	collection *gocb.Collection,
	cluster *gocb.Cluster,
) UserRepository {
	return &couchbaseUserRepository{
		logger:     logger.WithCodeNamespace(fmt.Sprintf("%T", &couchbaseUserRepository{})),
		tracer:     tracer,
		collection: collection,
		cluster:    cluster,
	}
}

func (repository *couchbaseUserRepository) Store(ctx context.Context, user *entities.User) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Insert(string(user.ID), user, &gocb.InsertOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot save user with ID [%s]", user.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseUserRepository) Update(ctx context.Context, user *entities.User) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Upsert(string(user.ID), user, &gocb.UpsertOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot update user with ID [%s]", user.ID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return nil
}

func (repository *couchbaseUserRepository) Load(ctx context.Context, userID entities.UserID) (*entities.User, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	result, err := repository.collection.Get(string(userID), &gocb.GetOptions{Context: ctx})
	if errors.Is(err, gocb.ErrDocumentNotFound) {
		msg := fmt.Sprintf("user with ID [%s] does not exist", userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}
	if err != nil {
		msg := fmt.Sprintf("cannot load user with ID [%s]", userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	user := new(entities.User)
	if err = result.Content(user); err != nil {
		msg := fmt.Sprintf("cannot decode user with ID [%s]", userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return user, nil
}

func (repository *couchbaseUserRepository) LoadBySubscriptionID(ctx context.Context, subscriptionID string) (*entities.User, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := fmt.Sprintf(
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.subscription_id = $subscriptionID",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: map[string]interface{}{"subscriptionID": subscriptionID},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot query user with subscriptionID [%s]", subscriptionID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repository.logger.Error(closeErr)
		}
	}()

	user := new(entities.User)
	if !rows.Next() {
		msg := fmt.Sprintf("user with subscriptionID [%s] does not exist", subscriptionID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	if err = rows.Row(user); err != nil {
		msg := fmt.Sprintf("cannot decode user with subscriptionID [%s]", subscriptionID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return user, nil
}

func (repository *couchbaseUserRepository) LoadOrStore(ctx context.Context, authUser entities.AuthUser) (*entities.User, bool, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	user := &entities.User{
		ID:               authUser.ID,
		Email:            authUser.Email,
		SubscriptionName: entities.SubscriptionNameFree,
		FirstName:        authUser.FirstName,
		LastName:         authUser.LastName,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	_, err := repository.collection.Insert(string(user.ID), user, &gocb.InsertOptions{Context: ctx})
	if err == nil {
		return user, true, nil
	}

	if !errors.Is(err, gocb.ErrDocumentExists) {
		msg := fmt.Sprintf("cannot load or create user from auth user [%+#v]", authUser)
		return nil, false, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	// Document already exists, load it
	existingUser, err := repository.Load(ctx, authUser.ID)
	if err != nil {
		msg := fmt.Sprintf("cannot load or create user from auth user [%+#v]", authUser)
		return nil, false, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return existingUser, false, nil
}
