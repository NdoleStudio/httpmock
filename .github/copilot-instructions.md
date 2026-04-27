# Copilot Instructions for httpmock

## Architecture

This is a monorepo with two main components:

- **`api/`** — Go backend (GoFiber v2) serving as the mock HTTP server and REST API. Uses PostgreSQL (via GORM), OpenTelemetry for observability, Clerk for auth, and LemonSqueezy for billing.
- **`web/`** — Next.js 16 frontend (React 19, TypeScript) using GitHub Primer design system, Clerk auth, and Zustand for state management.

The API follows a layered architecture: `handlers → validators → services → repositories`. Domain events flow through an event dispatcher with listeners (CloudEvents format via `cloudevents/sdk-go`). The DI container in `api/pkg/di/container.go` wires everything together.

API models are generated from Swagger into the frontend via `swagger-typescript-api` (see `web/package.json` script `api:model`).

## Build & Run

### API (Go)

```bash
cd api
go build -o main.exe .        # Build
go run .                       # Run (loads .env automatically)
swag init                      # Regenerate Swagger docs
```

### Web (Next.js)

```bash
cd web
pnpm install                   # Install dependencies
pnpm dev                       # Dev server with Turbopack
pnpm build                     # Production build
pnpm lint                      # ESLint
pnpm api:model                 # Regenerate TypeScript API models from Swagger
```

## Code Conventions

### Go (API)

- Error wrapping uses `github.com/palantir/stacktrace` (`stacktrace.Propagate`)
- Logging/tracing via custom `telemetry.Logger` and `telemetry.Tracer` interfaces
- Handlers embed a base `handler` struct for shared response helpers (`responseOK`, `responseBadRequest`, etc.)
- Repositories use interface + GORM implementation pattern (e.g., `ProjectRepository` interface, `gorm_project_repository.go` impl)
- IDs: UUIDs for entities (`github.com/google/uuid`), ULIDs for request logs (`github.com/oklog/ulid/v2`)
- Config loaded from environment variables via `github.com/caarlos0/env/v11`
- Routes registered in handler `RegisterRoutes` methods under versioned paths (`/v1/...`)

### TypeScript (Web)

- UI components from `@primer/react` (GitHub Primer design system, dark theme)
- State management via Zustand vanilla stores with a React context provider pattern (`src/store/`)
- API client is a configured Axios instance; types auto-generated from Swagger (`src/api/model.ts`)
- Toast notifications via `sonner`

### Formatting & Linting

Pre-commit hooks enforce: `go-fumpt`, `go-imports`, `go-lint`, `go-mod-tidy`, `prettier`, trailing whitespace removal, and YAML validation.
