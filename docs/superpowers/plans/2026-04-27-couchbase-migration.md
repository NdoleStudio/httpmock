# Couchbase Migration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the GORM/PostgreSQL persistence layer with Couchbase using `gocb/v2`, deleting all GORM code.

**Architecture:** Swap repository implementations behind existing interfaces. DI container provides `*gocb.Cluster` and per-entity `*gocb.Collection` instances. Indexes created on startup.

**Tech Stack:** Go 1.25, `github.com/couchbase/gocb/v2`, Couchbase Server, GoFiber v2

---

## File Structure

| Action | File                                                                | Responsibility                                            |
| ------ | ------------------------------------------------------------------- | --------------------------------------------------------- |
| Create | `pkg/repositories/couchbase_project_repository.go`                  | ProjectRepository implementation                          |
| Create | `pkg/repositories/couchbase_project_endpoint_repository.go`         | ProjectEndpointRepository implementation                  |
| Create | `pkg/repositories/couchbase_project_endpoint_request_repository.go` | ProjectEndpointRequestRepository implementation           |
| Create | `pkg/repositories/couchbase_user_repository.go`                     | UserRepository implementation                             |
| Modify | `pkg/di/container.go`                                               | Replace DB/GORM wiring with Couchbase cluster/collections |
| Modify | `pkg/entities/project.go`                                           | Remove `gorm:` struct tags                                |
| Modify | `pkg/entities/project_endpoint.go`                                  | Remove `gorm:` struct tags                                |
| Modify | `pkg/entities/project_endpoint_request.go`                          | Remove `gorm:` struct tags                                |
| Modify | `pkg/entities/user.go`                                              | Remove `gorm:` struct tags                                |
| Modify | `go.mod`                                                            | Add gocb/v2, remove GORM deps                             |
| Modify | `docker-compose.yml`                                                | Replace postgres with couchbase                           |
| Modify | `.env`                                                              | Replace DATABASE_URL with Couchbase vars                  |
| Delete | `pkg/repositories/gorm_project_repository.go`                       |                                                           |
| Delete | `pkg/repositories/gorm_project_endpoint_repository.go`              |                                                           |
| Delete | `pkg/repositories/gorm_project_endpoint_request_repository.go`      |                                                           |
| Delete | `pkg/repositories/gorm_user_repository.go`                          |                                                           |
| Delete | `pkg/telemetry/gorm_logger.go`                                      |                                                           |

---

### Task 1: Add Couchbase SDK dependency

**Files:**

- Modify: `api/go.mod`

- [ ] **Step 1: Add gocb/v2 to go.mod**

```bash
cd api
go get github.com/couchbase/gocb/v2@latest
```

- [ ] **Step 2: Remove GORM dependencies**

```bash
cd api
go get -u gorm.io/gorm@none gorm.io/driver/postgres@none gorm.io/plugin/opentelemetry@none github.com/dgraph-io/ristretto/v2@none
```

Note: This step may fail until we remove the code that imports them. We'll do a final `go mod tidy` after all code changes. For now, just add gocb/v2.

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "feat: add couchbase gocb/v2 SDK dependency

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 2: Remove GORM struct tags from entities

**Files:**

- Modify: `api/pkg/entities/project.go`
- Modify: `api/pkg/entities/project_endpoint.go`
- Modify: `api/pkg/entities/project_endpoint_request.go`
- Modify: `api/pkg/entities/user.go`

- [ ] **Step 1: Update `project.go`**

```go
package entities

import (
	"time"

	"github.com/google/uuid"
)

// Project is a  project belonging to a user
type Project struct {
	ID          uuid.UUID `json:"id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID      UserID    `json:"user_id" example:"user_2oeyIzOf9xxxxxxxxxxxxxx"`
	Subdomain   string    `json:"subdomain" example:"stripe-mock-api"`
	Name        string    `json:"name" example:"Mock Stripe API"`
	Description string    `json:"description" example:"Mock API for an online store for selling shoes"`
	CreatedAt   time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt   time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
```

- [ ] **Step 2: Update `project_endpoint.go`**

```go
package entities

import (
	"time"

	"github.com/google/uuid"
)

// ProjectEndpoint is an endpoint belonging to a project
type ProjectEndpoint struct {
	ID                          uuid.UUID `json:"id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	ProjectID                   uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	ProjectSubdomain            string    `json:"project_subdomain" example:"stripe-mock-api"`
	UserID                      UserID    `json:"user_id" example:"user_2oeyIzOf9xxxxxxxxxxxxxx"`
	RequestMethod               string    `json:"request_method" example:"GET"`
	RequestPath                 string    `json:"request_path" example:"/v1/products"`
	ResponseCode                uint      `json:"response_code" example:"200"`
	ResponseBody                *string   `json:"response_body" example:"{\"message\": \"Hello World\",\"status\": 200}"`
	ResponseHeaders             *string   `json:"response_headers" example:"[{\"Content-Type\":\"application/json\"}]"`
	ResponseDelayInMilliseconds uint      `json:"response_delay_in_milliseconds" example:"100"`
	Description                 *string   `json:"description" example:"Mock API for an online store for the /v1/products endpoint"`
	RequestCount                uint      `json:"request_count" example:"100"`
	CreatedAt                   time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt                   time.Time `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
```

- [ ] **Step 3: Update `project_endpoint_request.go`**

```go
package entities

import (
	"time"

	"github.com/google/uuid"
)

// ProjectEndpointRequest is the model for a project endpoint request
type ProjectEndpointRequest struct {
	ID                          string    `json:"id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	ProjectID                   uuid.UUID `json:"project_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	ProjectEndpointID           uuid.UUID `json:"project_endpoint_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	UserID                      UserID    `json:"user_id" example:"user_2oeyIzOf9xxxxxxxxxxxxxx"`
	RequestMethod               string    `json:"request_method" example:"GET"`
	RequestURL                  string    `json:"request_url" example:"https://stripe-mock-api.httpmock.dev/v1/products"`
	RequestHeaders              *string   `json:"request_headers" example:"[{\"Authorization\":\"Bearer sk_test_4eC39HqLyjWDarjtT1zdp7dc\"}]"`
	RequestBody                 *string   `json:"request_body" example:"{\"name\": \"Product 1\"}"`
	RequestIPAddress            string    `json:"request_ip_address" example:"127.0.0.1"`
	ResponseCode                uint      `json:"response_code" example:"200"`
	ResponseBody                *string   `json:"response_body" example:"{\"message\": \"Hello World\",\"status\": 200}"`
	ResponseHeaders             *string   `json:"response_headers" example:"[{\"Content-Type\":\"application/json\"}]"`
	ResponseDelayInMilliseconds uint      `json:"response_delay_in_milliseconds" example:"1000"`
	CreatedAt                   time.Time `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
}
```

- [ ] **Step 4: Update `user.go`**

```go
package entities

import (
	"time"
)

// SubscriptionName is the name of the subscription
type SubscriptionName string

// SubscriptionNameFree represents a free subscription
const SubscriptionNameFree = SubscriptionName("free")

// SubscriptionName10kMonthly represents a 10k pro subscription
const SubscriptionName10kMonthly = SubscriptionName("10k-monthly")

// SubscriptionName10kYearly represents a yearly pro subscription
const SubscriptionName10kYearly = SubscriptionName("100k-yearly")

// User stores information about a user
type User struct {
	ID                   UserID           `json:"id" example:"user_2oeyIzOf9xxxxxxxxxxxxxx"`
	Email                string           `json:"email" example:"name@email.com"`
	FirstName            *string          `json:"first_name" example:"John"`
	LastName             *string          `json:"last_name" example:"Doe"`
	SubscriptionName     SubscriptionName `json:"subscription_name" example:"free"`
	SubscriptionID       string           `json:"subscription_id" example:"8f9c71b8-b84e-4417-8408-a62274f65a08"`
	SubscriptionStatus   string           `json:"subscription_status" example:"free"`
	SubscriptionRenewsAt *time.Time       `json:"subscription_renews_at" example:"2022-06-05T14:26:02.302718+03:00"`
	SubscriptionEndsAt   *time.Time       `json:"subscription_ends_at" example:"2022-06-05T14:26:02.302718+03:00"`
	CreatedAt            time.Time        `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt            time.Time        `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
}
```

- [ ] **Step 5: Commit**

```bash
git add pkg/entities/
git commit -m "refactor: remove GORM struct tags from entities

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 3: Implement Couchbase User Repository

**Files:**

- Create: `api/pkg/repositories/couchbase_user_repository.go`

- [ ] **Step 1: Create the file**

```go
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

type couchbaseUserRepository struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

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
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.subscription_id = $subscriptionID LIMIT 1",
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

	user := new(entities.User)
	if !rows.Next() {
		msg := fmt.Sprintf("user with subscriptionID [%s] does not exist", subscriptionID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	if err = rows.Row(user); err != nil {
		msg := fmt.Sprintf("cannot decode user with subscriptionID [%s]", subscriptionID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
	}

	return user, nil
}

func (repository *couchbaseUserRepository) LoadOrStore(ctx context.Context, authUser entities.AuthUser) (user *entities.User, created bool, err error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	user = &entities.User{
		ID:               authUser.ID,
		Email:            authUser.Email,
		SubscriptionName: entities.SubscriptionNameFree,
		FirstName:        authUser.FirstName,
		LastName:         authUser.LastName,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	_, err = repository.collection.Insert(string(user.ID), user, &gocb.InsertOptions{Context: ctx})
	if err == nil {
		return user, true, nil
	}

	if !errors.Is(err, gocb.ErrDocumentExists) {
		msg := fmt.Sprintf("cannot insert user from auth user [%+#v]", authUser)
		return nil, false, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	// Document exists, load it
	existingUser, err := repository.Load(ctx, authUser.ID)
	if err != nil {
		msg := fmt.Sprintf("cannot load existing user with ID [%s]", authUser.ID)
		return nil, false, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	return existingUser, false, nil
}
```

- [ ] **Step 2: Commit**

```bash
git add pkg/repositories/couchbase_user_repository.go
git commit -m "feat: add Couchbase user repository implementation

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 4: Implement Couchbase Project Repository

**Files:**

- Create: `api/pkg/repositories/couchbase_project_repository.go`

- [ ] **Step 1: Create the file**

```go
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

type couchbaseProjectRepository struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

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

	var projects []*entities.Project
	for rows.Next() {
		project := new(entities.Project)
		if err = rows.Row(project); err != nil {
			msg := fmt.Sprintf("cannot decode project for user [%s]", userID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		projects = append(projects, project)
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
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

	// Verify ownership
	if project.UserID != userID {
		msg := fmt.Sprintf("project with ID [%s] does not belong to user [%s]", projectID, userID)
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
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.subdomain = $subdomain LIMIT 1",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: map[string]interface{}{"subdomain": subdomain},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot query project with subdomain [%s]", subdomain)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	project := new(entities.Project)
	if !rows.Next() {
		msg := fmt.Sprintf("project not found with subdomain [%s]", subdomain)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	if err = rows.Row(project); err != nil {
		msg := fmt.Sprintf("cannot decode project with subdomain [%s]", subdomain)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
	}

	return project, nil
}
```

- [ ] **Step 2: Commit**

```bash
git add pkg/repositories/couchbase_project_repository.go
git commit -m "feat: add Couchbase project repository implementation

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 5: Implement Couchbase Project Endpoint Repository

**Files:**

- Create: `api/pkg/repositories/couchbase_project_endpoint_repository.go`

- [ ] **Step 1: Create the file**

```go
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

type couchbaseProjectEndpointRepository struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

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
		msg := fmt.Sprintf("cannot increase request_count for endpoint with ID [%s]", projectEndpointID)
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
		msg := fmt.Sprintf("cannot decrease request_count for endpoint with ID [%s]", projectEndpointID)
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
		msg := fmt.Sprintf("cannot update project_subdomain for project ID [%s]", projectID)
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
		msg := fmt.Sprintf("cannot load project endpoints for user [%s] and project [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	var endpoints []*entities.ProjectEndpoint
	for rows.Next() {
		endpoint := new(entities.ProjectEndpoint)
		if err = rows.Row(endpoint); err != nil {
			msg := fmt.Sprintf("cannot decode project endpoint for user [%s]", userID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		endpoints = append(endpoints, endpoint)
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
	}

	return endpoints, nil
}

func (repository *couchbaseProjectEndpointRepository) Load(ctx context.Context, userID entities.UserID, projectID uuid.UUID, projectEndpointID uuid.UUID) (*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	result, err := repository.collection.Get(projectEndpointID.String(), &gocb.GetOptions{Context: ctx})
	if errors.Is(err, gocb.ErrDocumentNotFound) {
		msg := fmt.Sprintf("endpoint with ID [%s] does not exist", projectEndpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg))
	}
	if err != nil {
		msg := fmt.Sprintf("cannot load endpoint with ID [%s]", projectEndpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	endpoint := new(entities.ProjectEndpoint)
	if err = result.Content(endpoint); err != nil {
		msg := fmt.Sprintf("cannot decode endpoint with ID [%s]", projectEndpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if endpoint.UserID != userID || endpoint.ProjectID != projectID {
		msg := fmt.Sprintf("endpoint with ID [%s] does not belong to user [%s] project [%s]", projectEndpointID, userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	return endpoint, nil
}

func (repository *couchbaseProjectEndpointRepository) Delete(ctx context.Context, endpoint *entities.ProjectEndpoint) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	_, err := repository.collection.Remove(endpoint.ID.String(), &gocb.RemoveOptions{Context: ctx})
	if err != nil {
		msg := fmt.Sprintf("cannot delete endpoint with ID [%s] for user [%s]", endpoint.ID, endpoint.UserID)
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	return nil
}

func (repository *couchbaseProjectEndpointRepository) LoadByRequestForUser(ctx context.Context, userID entities.UserID, projectID uuid.UUID, requestMethod, requestPath string) (*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := fmt.Sprintf(
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_id = $projectID AND d.request_path = $requestPath AND (d.request_method = $requestMethod OR d.request_method = 'ANY') LIMIT 1",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	params := map[string]interface{}{
		"userID":        string(userID),
		"projectID":    projectID.String(),
		"requestPath":  requestPath,
		"requestMethod": strings.ToUpper(requestMethod),
	}

	if requestMethod == "ANY" {
		query = fmt.Sprintf(
			"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_id = $projectID AND d.request_path = $requestPath LIMIT 1",
			repository.collection.Bucket().Name(),
			repository.collection.ScopeName(),
			repository.collection.Name(),
		)
	}

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: params,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot query endpoint with method [%s] path [%s]", requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	endpoint := new(entities.ProjectEndpoint)
	if !rows.Next() {
		msg := fmt.Sprintf("endpoint not found with project ID [%s], request method [%s] and request path [%s]", projectID, requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	if err = rows.Row(endpoint); err != nil {
		msg := fmt.Sprintf("cannot decode endpoint for project [%s]", projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
	}

	return endpoint, nil
}

func (repository *couchbaseProjectEndpointRepository) LoadByRequest(ctx context.Context, subdomain, requestMethod, requestPath string) (*entities.ProjectEndpoint, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := fmt.Sprintf(
		"SELECT d.* FROM `%s`.`%s`.`%s` d WHERE d.project_subdomain = $subdomain AND d.request_path = $requestPath AND (d.request_method = $requestMethod OR d.request_method = 'ANY') LIMIT 1",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
		NamedParameters: map[string]interface{}{
			"subdomain":     subdomain,
			"requestMethod": strings.ToUpper(requestMethod),
			"requestPath":   requestPath,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot query endpoint with subdomain [%s] method [%s] path [%s]", subdomain, requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	endpoint := new(entities.ProjectEndpoint)
	if !rows.Next() {
		msg := fmt.Sprintf("endpoint not found with subdomain [%s] request method [%s] and request path [%s]", subdomain, requestMethod, requestPath)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	if err = rows.Row(endpoint); err != nil {
		msg := fmt.Sprintf("cannot decode endpoint for subdomain [%s]", subdomain)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
	}

	return endpoint, nil
}
```

- [ ] **Step 2: Commit**

```bash
git add pkg/repositories/couchbase_project_endpoint_repository.go
git commit -m "feat: add Couchbase project endpoint repository implementation

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 6: Implement Couchbase Project Endpoint Request Repository

**Files:**

- Create: `api/pkg/repositories/couchbase_project_endpoint_request_repository.go`

- [ ] **Step 1: Create the file**

```go
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

type couchbaseProjectEndpointRequestRepository struct {
	logger     telemetry.Logger
	tracer     telemetry.Tracer
	collection *gocb.Collection
	cluster    *gocb.Cluster
}

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
		msg := fmt.Sprintf("cannot delete request with ID [%s] for user [%s]", request.ID, request.UserID)
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
		msg := fmt.Sprintf("cannot load request with ID [%s]", requestID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	request := new(entities.ProjectEndpointRequest)
	if err = result.Content(request); err != nil {
		msg := fmt.Sprintf("cannot decode request with ID [%s]", requestID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	if request.UserID != userID {
		msg := fmt.Sprintf("request with ID [%s] does not belong to user [%s]", requestID, userID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCode(ErrCodeNotFound, msg))
	}

	return request, nil
}

func (repository *couchbaseProjectEndpointRequestRepository) GetProjectTraffic(ctx context.Context, userID entities.UserID, projectID uuid.UUID) ([]*TimeSeriesData, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	thirtyDaysAgo := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	query := fmt.Sprintf(
		"SELECT DATE_TRUNC_STR(d.created_at, 'day') AS timestamp, COUNT(*) AS `count` FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_id = $projectID AND d.created_at >= $thirtyDaysAgo GROUP BY DATE_TRUNC_STR(d.created_at, 'day')",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
		NamedParameters: map[string]interface{}{
			"userID":        string(userID),
			"projectID":    projectID.String(),
			"thirtyDaysAgo": thirtyDaysAgo,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load project traffic for user [%s] and project [%s]", userID, projectID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	var data []*TimeSeriesData
	for rows.Next() {
		point := new(TimeSeriesData)
		if err = rows.Row(point); err != nil {
			msg := fmt.Sprintf("cannot decode traffic data for project [%s]", projectID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		data = append(data, point)
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
	}

	return repository.normalizeTimeSeries(data), nil
}

func (repository *couchbaseProjectEndpointRequestRepository) GetEndpointTraffic(ctx context.Context, userID entities.UserID, endpointID uuid.UUID) ([]*TimeSeriesData, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	thirtyDaysAgo := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)

	query := fmt.Sprintf(
		"SELECT DATE_TRUNC_STR(d.created_at, 'day') AS timestamp, COUNT(*) AS `count` FROM `%s`.`%s`.`%s` d WHERE d.user_id = $userID AND d.project_endpoint_id = $endpointID AND d.created_at >= $thirtyDaysAgo GROUP BY DATE_TRUNC_STR(d.created_at, 'day')",
		repository.collection.Bucket().Name(),
		repository.collection.ScopeName(),
		repository.collection.Name(),
	)

	rows, err := repository.cluster.Query(query, &gocb.QueryOptions{
		Context: ctx,
		NamedParameters: map[string]interface{}{
			"userID":        string(userID),
			"endpointID":   endpointID.String(),
			"thirtyDaysAgo": thirtyDaysAgo,
		},
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load endpoint traffic for user [%s] and endpoint [%s]", userID, endpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	var data []*TimeSeriesData
	for rows.Next() {
		point := new(TimeSeriesData)
		if err = rows.Row(point); err != nil {
			msg := fmt.Sprintf("cannot decode traffic data for endpoint [%s]", endpointID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		data = append(data, point)
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
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

	if previousID != nil {
		baseQuery += " AND d.id < $cursorID ORDER BY d.id DESC LIMIT $limit"
		params["cursorID"] = previousID.String()
	} else if nextID != nil {
		baseQuery += " AND d.id > $cursorID ORDER BY d.id ASC LIMIT $limit"
		params["cursorID"] = nextID.String()
	} else {
		baseQuery += " ORDER BY d.id DESC LIMIT $limit"
	}

	rows, err := repository.cluster.Query(baseQuery, &gocb.QueryOptions{
		Context:         ctx,
		NamedParameters: params,
	})
	if err != nil {
		msg := fmt.Sprintf("cannot load requests for user [%s] and endpoint [%s]", userID, endpointID)
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}

	var requests []*entities.ProjectEndpointRequest
	for rows.Next() {
		request := new(entities.ProjectEndpointRequest)
		if err = rows.Row(request); err != nil {
			msg := fmt.Sprintf("cannot decode request for endpoint [%s]", endpointID)
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
		}
		requests = append(requests, request)
	}

	if err = rows.Close(); err != nil {
		repository.logger.Error(stacktrace.Propagate(err, "cannot close query rows"))
	}

	return requests, nil
}

func (repository *couchbaseProjectEndpointRequestRepository) generateTimeSeries() map[string]*TimeSeriesData {
	series := make(map[string]*TimeSeriesData)
	for i := 0; i < 30; i++ {
		date := time.Now().UTC().AddDate(0, 0, -i)
		series[date.Format("2006-01-02")] = &TimeSeriesData{
			Timestamp: time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
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
```

- [ ] **Step 2: Commit**

```bash
git add pkg/repositories/couchbase_project_endpoint_request_repository.go
git commit -m "feat: add Couchbase project endpoint request repository implementation

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 7: Rewrite DI container for Couchbase

**Files:**

- Modify: `api/pkg/di/container.go`

- [ ] **Step 1: Replace the Container struct and imports**

Remove GORM/ristretto imports. Add `gocb/v2` import. Replace `db *gorm.DB` field with `cluster *gocb.Cluster` and `bucket *gocb.Bucket`.

Updated struct:

```go
type Container struct {
	projectID       string
	version         string
	cluster         *gocb.Cluster
	bucket          *gocb.Bucket
	app             *fiber.App
	eventDispatcher *services.EventDispatcher
	logger          telemetry.Logger
}
```

- [ ] **Step 2: Add `Cluster()` method**

```go
func (container *Container) Cluster() *gocb.Cluster {
	if container.cluster != nil {
		return container.cluster
	}

	container.logger.Debug("creating *gocb.Cluster")

	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_CONNECTION_STRING"), gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: os.Getenv("COUCHBASE_USERNAME"),
			Password: os.Getenv("COUCHBASE_PASSWORD"),
		},
		Tracer: gocb.NewOpenTelemetryRequestTracer(container.Tracer().Provider()),
	})
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot connect to Couchbase cluster"))
	}

	err = cluster.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "Couchbase cluster not ready"))
	}

	container.cluster = cluster
	return cluster
}
```

- [ ] **Step 3: Add `Bucket()` and collection methods**

```go
func (container *Container) Bucket() *gocb.Bucket {
	if container.bucket != nil {
		return container.bucket
	}

	container.logger.Debug("creating *gocb.Bucket")
	container.bucket = container.Cluster().Bucket(os.Getenv("COUCHBASE_BUCKET"))
	err := container.bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "Couchbase bucket not ready"))
	}
	return container.bucket
}

func (container *Container) ProjectsCollection() *gocb.Collection {
	return container.Bucket().Scope("_default").Collection("projects")
}

func (container *Container) EndpointsCollection() *gocb.Collection {
	return container.Bucket().Scope("_default").Collection("project_endpoints")
}

func (container *Container) EndpointRequestsCollection() *gocb.Collection {
	return container.Bucket().Scope("_default").Collection("project_endpoint_requests")
}

func (container *Container) UsersCollection() *gocb.Collection {
	return container.Bucket().Scope("_default").Collection("users")
}
```

- [ ] **Step 4: Add `EnsureCollections()` method**

```go
func (container *Container) EnsureCollections() {
	container.logger.Debug("ensuring Couchbase collections exist")
	collections := container.Bucket().Collections()

	collectionNames := []string{"projects", "project_endpoints", "project_endpoint_requests", "users"}
	for _, name := range collectionNames {
		err := collections.CreateCollection(gocb.CollectionSpec{
			Name:      name,
			ScopeName: "_default",
		}, nil)
		if err != nil && !errors.Is(err, gocb.ErrCollectionExists) {
			container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot create collection [%s]", name)))
		}
	}
}
```

- [ ] **Step 5: Add `EnsureIndexes()` method**

```go
func (container *Container) EnsureIndexes() {
	container.logger.Debug("ensuring Couchbase indexes exist")
	bucket := os.Getenv("COUCHBASE_BUCKET")

	indexes := []string{
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_projects_user_id ON `%s`.`_default`.`projects`(user_id)", bucket),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_projects_subdomain ON `%s`.`_default`.`projects`(subdomain)", bucket),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_endpoints_user_project ON `%s`.`_default`.`project_endpoints`(user_id, project_id)", bucket),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_endpoints_subdomain_request ON `%s`.`_default`.`project_endpoints`(project_subdomain, request_method, request_path)", bucket),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_endpoints_project_id ON `%s`.`_default`.`project_endpoints`(project_id)", bucket),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_requests_user_endpoint ON `%s`.`_default`.`project_endpoint_requests`(user_id, project_endpoint_id, id DESC)", bucket),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_requests_user_project_created ON `%s`.`_default`.`project_endpoint_requests`(user_id, project_id, created_at)", bucket),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_requests_user_endpoint_created ON `%s`.`_default`.`project_endpoint_requests`(user_id, project_endpoint_id, created_at)", bucket),
		fmt.Sprintf("CREATE INDEX IF NOT EXISTS idx_users_subscription_id ON `%s`.`_default`.`users`(subscription_id)", bucket),
	}

	for _, query := range indexes {
		_, err := container.Cluster().Query(query, nil)
		if err != nil {
			container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot create index: %s", query)))
		}
	}
}
```

- [ ] **Step 6: Update repository factory methods**

Replace the 4 repository methods:

```go
func (container *Container) ProjectRepository() repositories.ProjectRepository {
	container.logger.Debug("creating Couchbase repositories.ProjectRepository")
	return repositories.NewCouchbaseProjectRepository(
		container.Logger(),
		container.Tracer(),
		container.ProjectsCollection(),
		container.Cluster(),
	)
}

func (container *Container) ProjectEndpointRepository() repositories.ProjectEndpointRepository {
	container.logger.Debug("creating Couchbase repositories.ProjectEndpointRepository")
	return repositories.NewCouchbaseProjectEndpointRepository(
		container.Logger(),
		container.Tracer(),
		container.EndpointsCollection(),
		container.Cluster(),
	)
}

func (container *Container) ProjectEndpointRequestRepository() repositories.ProjectEndpointRequestRepository {
	container.logger.Debug("creating Couchbase repositories.ProjectEndpointRequestRepository")
	return repositories.NewCouchbaseProjectEndpointRequestRepository(
		container.Logger(),
		container.Tracer(),
		container.EndpointRequestsCollection(),
		container.Cluster(),
	)
}

func (container *Container) UserRepository() repositories.UserRepository {
	container.logger.Debug("creating Couchbase repositories.UserRepository")
	return repositories.NewCouchbaseUserRepository(
		container.Logger(),
		container.Tracer(),
		container.UsersCollection(),
		container.Cluster(),
	)
}
```

- [ ] **Step 7: Update `App()` method**

Add `container.EnsureCollections()` and `container.EnsureIndexes()` calls at the start of `App()`, before registering routes.

- [ ] **Step 8: Remove dead code**

Delete the `DB()`, `GormLogger()`, and `ProjectEndpointRistrettoCache()` methods entirely.

- [ ] **Step 9: Remove unused imports**

Remove these imports from container.go:

- `"gorm.io/gorm"`
- `"gorm.io/driver/postgres"`
- `gormLogger "gorm.io/gorm/logger"`
- `"gorm.io/plugin/opentelemetry/tracing"`
- `"github.com/dgraph-io/ristretto/v2"`

Add:

- `"github.com/couchbase/gocb/v2"`

- [ ] **Step 10: Commit**

```bash
git add pkg/di/container.go
git commit -m "feat: rewire DI container from GORM to Couchbase

- Add Cluster(), Bucket(), collection accessors
- Add EnsureCollections() and EnsureIndexes() for startup provisioning
- Replace repository factory methods to use Couchbase implementations
- Remove DB(), GormLogger(), ProjectEndpointRistrettoCache()

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 8: Delete GORM repository files and telemetry

**Files:**

- Delete: `api/pkg/repositories/gorm_project_repository.go`
- Delete: `api/pkg/repositories/gorm_project_endpoint_repository.go`
- Delete: `api/pkg/repositories/gorm_project_endpoint_request_repository.go`
- Delete: `api/pkg/repositories/gorm_user_repository.go`
- Delete: `api/pkg/telemetry/gorm_logger.go`

- [ ] **Step 1: Delete the files**

```bash
cd api
rm pkg/repositories/gorm_project_repository.go
rm pkg/repositories/gorm_project_endpoint_repository.go
rm pkg/repositories/gorm_project_endpoint_request_repository.go
rm pkg/repositories/gorm_user_repository.go
rm pkg/telemetry/gorm_logger.go
```

- [ ] **Step 2: Commit**

```bash
git add -A
git commit -m "refactor: delete GORM repository implementations and logger

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 9: Update go.mod and tidy dependencies

**Files:**

- Modify: `api/go.mod`
- Modify: `api/go.sum`

- [ ] **Step 1: Run go mod tidy**

```bash
cd api
go mod tidy
```

- [ ] **Step 2: Verify build compiles**

```bash
cd api
go build ./...
```

If there are compile errors, fix them. Common issues:

- Unused imports remaining in container.go
- Missing `"errors"` import in container.go (needed for `errors.Is`)

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: tidy go.mod, remove GORM dependencies

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 10: Update infrastructure configuration

**Files:**

- Modify: `api/docker-compose.yml`
- Modify: `api/.env`

- [ ] **Step 1: Update docker-compose.yml**

Replace the entire file with:

```yaml
services:
  api:
    image: ndolestudio/httpmock:latest
    ports:
      - "80:8000"
      - "443:8443"
    env_file: .env
    volumes:
      - ./certs:/app/certs
    depends_on:
      - couchbase

  couchbase:
    image: couchbase/server:latest
    ports:
      - "8091:8091"
      - "8093:8093"
      - "11210:11210"
    environment:
      - COUCHBASE_ADMINISTRATOR_USERNAME=Administrator
      - COUCHBASE_ADMINISTRATOR_PASSWORD=password
```

- [ ] **Step 2: Update .env**

Replace `DATABASE_URL=...` line with:

```
COUCHBASE_CONNECTION_STRING=couchbase://localhost
COUCHBASE_USERNAME=Administrator
COUCHBASE_PASSWORD=password
COUCHBASE_BUCKET=httpmock
```

Also remove `DB_USERNAME` and `DB_PASSWORD` lines.

- [ ] **Step 3: Update Dockerfile**

No changes needed — it just builds the Go binary.

- [ ] **Step 4: Commit**

```bash
git add docker-compose.yml .env
git commit -m "feat: update infrastructure for Couchbase

- Replace PostgreSQL with Couchbase in docker-compose
- Update .env with Couchbase connection variables

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```

---

### Task 11: Final verification build

- [ ] **Step 1: Full build check**

```bash
cd api
go build ./...
```

- [ ] **Step 2: Verify no GORM references remain**

```bash
grep -r "gorm" pkg/ --include="*.go" | grep -v "_test.go"
```

Expected: no output

- [ ] **Step 3: Verify no ristretto references remain**

```bash
grep -r "ristretto" pkg/ --include="*.go"
```

Expected: no output

- [ ] **Step 4: Final commit if any fixups needed**

```bash
git add -A
git commit -m "fix: resolve any remaining build issues

Co-authored-by: Copilot <223556219+Copilot@users.noreply.github.com>"
```
