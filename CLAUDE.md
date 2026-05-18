# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

**simpengine** — Go backend for [simpengine.xyz](https://simpengine.xyz). Users build multi-step "romantic events" (drinks → entertainment → food → ...) for a *simp target*, publish a public token link, and the target picks options. Stack: Go 1.24, Postgres 17, MinIO (object storage), JWT auth via cookie. The Next.js client lives in the sibling repo `simpengine-nextjs/`.

## Commands

Local development is Docker-first:

```bash
docker-compose up --build           # backend :8080, postgres :5433, minio :9000/:9001
```

Migrations run automatically on container start via `entrypoint.sh` (uses [golang-migrate](https://github.com/golang-migrate/migrate)). To run them manually against the dockerized DB:

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/simpengine?sslmode=disable" up
migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/simpengine?sslmode=disable" down 1
```

To create a new migration (sequence numbers are zero-padded 6 digits — see existing files like `000010_*`):

```bash
migrate create -ext sql -dir migrations -seq <name>
```

Build / run without Docker:

```bash
go build -o main ./cmd/web
go run ./cmd/web                    # needs .env or shell env: ADDR, DSN, JWT_SECRET, MINIO_*, FRONTEND_HOST, TELEGRAM_BOT_TOKEN
```

Regenerate Swagger docs (annotations live as `// @...` comments in handlers/`main.go`; output written to `cmd/web/docs/`):

```bash
swag init -g cmd/web/main.go -o cmd/web/docs
```

Swagger UI is served at `/swagger/` when the server runs.

## Architecture

### Module layout (`internal/<module>/`)

Each business domain is a self-contained module that follows the same pattern:

- `routes.go` — defines `Module` struct, `NewModule(deps...)` constructor that wires repo→service, and `RegisterRoutes(mux, cfg)` that attaches handlers
- `handler.go` — HTTP handlers with Swagger annotations
- `service.go` — business logic; takes interface-typed repositories
- `repository.go` (interface) + `postgres_repo.go` *or* `repository/postgres/*.go` (impl)
- `model.go` / `model/` — domain types
- `dto.go` — request/response shapes with `validator` tags

Modules: `user`, `auth`, `simp-target`, `romantic_event`, `media`, `notification`, `admin_invite`. Cross-module composition happens in `cmd/web/routes.go`: `romanticevent.NewModule` is constructed with the *services* of `simptarget`, `media`, `user`, and `notification` — modules talk to each other only through service interfaces, never repositories.

### Entry point and wiring

`cmd/web/main.go` reads env (loading `.env` if present), opens the Postgres pool, creates the MinIO client, builds `config.Config`, and starts the server. `routes()` (in `cmd/web/routes.go`) instantiates every module and chains middleware via `alice`: `recoverPanic → logRequest → secureHeaders`. The CORS allow-origin is `cfg.FrontEndHost` (single origin, with credentials).

### Auth model

JWT lives in an `HttpOnly`, `Secure`, `SameSite=None` cookie named `jwt`, set by `auth.Login`. `auth.AuthMiddleware` parses it and stuffs `user_id` (int64) + `role` (string) into the request context under unexported keys; downstream code reads them via `auth.GetUserIDFromContext` / `auth.GetRoleFromContext`. `auth.RoleMiddleware("admin")` wraps handlers that require a specific role (used for `/admin/template-event/*` routes).

### Error handling contract

All errors flow through `internal/shared/model/error.go`. Handlers call `model.WriteErrorResponse(w, err)`; the helper inspects sentinel errors (`ErrNoRecord`, `ErrInvalidBody`, `ErrInvalidCredentials`, `ErrForbidden`, `ErrUniqueViolation`, `ErrInvalidState`, `ErrValidation`, `ErrInvalidParams`) via `errors.Is` and maps them to HTTP status + machine code + message. Services wrap underlying errors with `fmt.Errorf("%w: ...", model.ErrX, ...)` so the sentinel survives.

Postgres errors get translated to these sentinels in `internal/shared/db/helpers.go::MapPQError` (unique violation → `ErrUniqueViolation`, FK violation → `ErrNoRecord`). Always pipe `db.MapPQError(err)` through repository methods rather than returning raw `*pq.Error`.

### Romantic event domain

The core entity has a state machine implied by `status` (`draft` / published / etc — see `internal/romantic_event/model/`). Mutations (update/delete steps/options) require the event to be in `draft`, enforced in the service layer with `ErrInvalidState`. Publishing assigns a `public_token` (UUID) used by the unauthenticated `/public/romantic-event/{public_token}` endpoints the target uses to view and submit choices.

### Media

`media` module has two repositories: Postgres (`media/repository/postgres`) for metadata rows and MinIO (`media/repository/minio`) for blobs. Both are wired together in `media.NewModule`. Uploaded files get a row in `media` and an object in the configured `MINIO_BUCKET`.

### Config / env vars

`config.Config` (in `cmd/config/config.go`) is the single struct passed to every `RegisterRoutes`. Env vars consumed in `main.go`: `ADDR`, `DSN`, `JWT_SECRET`, `MINIO_ENDPOINT`, `MINIO_BUCKET`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `FRONTEND_HOST`, `TELEGRAM_BOT_TOKEN`. Defaults in `docker-compose.yml` are for local dev only.
