# Database Migration: PostgreSQL/GORM → Couchbase

## Summary

Replace the GORM/PostgreSQL persistence layer with Couchbase using the `gocb/v2` SDK. The existing repository interfaces remain unchanged — only the implementations, DI wiring, and infrastructure configuration change.

**References:**

- https://docs.couchbase.com/go-sdk/current/hello-world/start-using-sdk.html
- https://pkg.go.dev/github.com/couchbase/gocb/v2

## Data Model

Single bucket `httpmock`, default scope, one collection per entity:

| Collection                  | Document Key                    | Entity                            |
| --------------------------- | ------------------------------- | --------------------------------- |
| `projects`                  | `project.ID.String()` (UUID)    | `entities.Project`                |
| `project_endpoints`         | `endpoint.ID.String()` (UUID)   | `entities.ProjectEndpoint`        |
| `project_endpoint_requests` | `request.ID` (ULID string)      | `entities.ProjectEndpointRequest` |
| `users`                     | `user.ID` (Clerk UserID string) | `entities.User`                   |

Documents are serialized as JSON using existing `json:"..."` struct tags. All `gorm:"..."` struct tags are removed from entities.

## Indexes

Created programmatically on startup via N1QL `CREATE INDEX IF NOT EXISTS`:

### projects

- `idx_projects_user_id` on `(user_id)` — `Fetch` by user
- `idx_projects_subdomain` on `(subdomain)` — `LoadWithSubdomain`

### project_endpoints

- `idx_endpoints_user_project` on `(user_id, project_id)` — `Fetch`, `Load`, `LoadByRequestForUser`
- `idx_endpoints_subdomain_request` on `(project_subdomain, request_method, request_path)` — `LoadByRequest` (hot path)
- `idx_endpoints_project_id` on `(project_id)` — `UpdateSubdomain`

### project_endpoint_requests

- `idx_requests_user_endpoint` on `(user_id, project_endpoint_id, id DESC)` — `Index` (paginated listing)
- `idx_requests_user_project_created` on `(user_id, project_id, created_at)` — `GetProjectTraffic`
- `idx_requests_user_endpoint_created` on `(user_id, project_endpoint_id, created_at)` — `GetEndpointTraffic`

### users

- `idx_users_subscription_id` on `(subscription_id)` — `LoadBySubscriptionID`

## Repository Implementations

### Files Created

- `couchbase_project_repository.go`
- `couchbase_project_endpoint_repository.go`
- `couchbase_project_endpoint_request_repository.go`
- `couchbase_user_repository.go`

### Files Deleted

- `gorm_project_repository.go`
- `gorm_project_endpoint_repository.go`
- `gorm_project_endpoint_request_repository.go`
- `gorm_user_repository.go`

### Constructor Pattern

Each repository receives its specific `*gocb.Collection`:

```go
func NewCouchbaseProjectRepository(
    logger telemetry.Logger,
    tracer telemetry.Tracer,
    collection *gocb.Collection,
) ProjectRepository
```

### Operation Mapping

| Operation                 | GORM                          | Couchbase                                                               |
| ------------------------- | ----------------------------- | ----------------------------------------------------------------------- |
| Create                    | `db.Create(entity)`           | `collection.Insert(id, entity, nil)`                                    |
| Read by ID                | `db.First(entity, id)`        | `collection.Get(id, nil)`                                               |
| Update                    | `db.Save(entity)`             | `collection.Upsert(id, entity, nil)`                                    |
| Delete                    | `db.Delete(entity)`           | `collection.Remove(id, nil)`                                            |
| Query/Filter              | `db.Where(...).Find(...)`     | `cluster.Query(n1ql, ...)`                                              |
| Atomic increment          | `gorm.Expr("field + 1")`      | `collection.MutateIn(id, []gocb.MutateInSpec{gocb.IncrementSpec(...)})` |
| Transaction (LoadOrStore) | `db.Transaction(func(tx)...)` | Try `Insert`; on `ErrDocumentExists`, fall back to `Get`                |

### Error Mapping

- `gocb.ErrDocumentNotFound` → `stacktrace.PropagateWithCode(err, ErrCodeNotFound, msg)`
- `gocb.ErrDocumentExists` → used in `LoadOrStore` for optimistic insert fallback

### Traffic Queries (N1QL)

```sql
SELECT DATE_TRUNC_STR(created_at, "day") AS timestamp, COUNT(*) AS count
FROM project_endpoint_requests
WHERE user_id = $userID AND project_id = $projectID
  AND created_at >= $thirtyDaysAgo
GROUP BY DATE_TRUNC_STR(created_at, "day")
```

The `normalizeTimeSeries` helper remains in Go (fills in zero-count days for the 30-day window).

## DI Container Changes

### New Methods

- `Cluster() *gocb.Cluster` — connects using `COUCHBASE_CONNECTION_STRING`, `COUCHBASE_USERNAME`, `COUCHBASE_PASSWORD`
- `Bucket() *gocb.Bucket` — opens `COUCHBASE_BUCKET`
- `ProjectsCollection() *gocb.Collection`
- `EndpointsCollection() *gocb.Collection`
- `EndpointRequestsCollection() *gocb.Collection`
- `UsersCollection() *gocb.Collection`
- `EnsureIndexes()` — runs `CREATE INDEX IF NOT EXISTS` statements

### Removed

- `DB() *gorm.DB`
- `GormLogger() gormLogger.Interface`
- Ristretto cache initialization
- `gorm.io/plugin/opentelemetry/tracing` usage

### OpenTelemetry Integration

Use `gocb.NewOpenTelemetryRequestTracer(otelTracer)` in `ClusterOptions` to maintain observability.

## Dependency Changes

### go.mod — Add

- `github.com/couchbase/gocb/v2`

### go.mod — Remove

- `gorm.io/gorm`
- `gorm.io/driver/postgres`
- `gorm.io/plugin/opentelemetry`
- `github.com/dgraph-io/ristretto/v2`

## Infrastructure Changes

### Environment Variables

- Remove: `DATABASE_URL`
- Add: `COUCHBASE_CONNECTION_STRING`, `COUCHBASE_USERNAME`, `COUCHBASE_PASSWORD`, `COUCHBASE_BUCKET`

### docker-compose.yml

Replace PostgreSQL container with:

```yaml
couchbase:
  image: couchbase/server:latest
  ports:
    - "8091:8091" # Web UI
    - "8093:8093" # N1QL
    - "11210:11210" # KV
  environment:
    - COUCHBASE_ADMINISTRATOR_USERNAME=Administrator
    - COUCHBASE_ADMINISTRATOR_PASSWORD=password
```

### .env

```
COUCHBASE_CONNECTION_STRING=couchbase://localhost
COUCHBASE_USERNAME=Administrator
COUCHBASE_PASSWORD=password
COUCHBASE_BUCKET=httpmock
```

## Scope Boundaries

- Repository interfaces: **unchanged**
- Services, handlers, validators, events, listeners: **unchanged**
- Frontend: **unchanged**
- Swagger docs: **unchanged**
