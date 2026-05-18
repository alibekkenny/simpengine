# Module structure

Every business domain lives under `internal/<module>/` and follows the same skeleton.

## Required files

```
internal/<module>/
├── routes.go         # Module struct, NewModule(deps), RegisterRoutes(mux, cfg)
├── handler.go        # HTTP handlers (+ Swagger annotations)
├── service.go        # Business logic; depends on repo interfaces
├── repository.go     # Repo interface(s)
├── postgres_repo.go  # Single-entity postgres impl (simple modules)
├── model.go          # Domain types (single-entity modules)
└── dto.go            # Request/response shapes with validator tags
```

## When a module has multiple entities

Split repositories and models into subpackages — e.g. `romantic_event` and `media`:

```
internal/<module>/
├── routes.go
├── handler.go
├── service.go
├── dto.go
├── model/
│   ├── <entity_a>.go
│   └── <entity_b>.go
└── repository/
    ├── <entity_a>.go            # interface
    ├── <entity_b>.go            # interface
    └── postgres/
        ├── <entity_a>.go        # impl
        └── <entity_b>.go        # impl
```

Reference: `internal/romantic_event/`, `internal/media/`.

## Package names

- Folder `simp-target` → `package simptarget` (drop the hyphen, lowercase).
- Folder `romantic_event` → `package romanticevent` for top-level, `package model` / `package repository` / `package postgres` for subpackages.
- Always import shared error sentinels as `"github.com/alibekkenny/simpengine/internal/shared/model"` and alias to `shared_model` when the module also defines its own `model` subpackage (see `internal/romantic_event/handler.go:10`).

## Naming

- Constructor: `NewModule(...)` returns `*Module`.
- Service constructor: `NewXService(repo)` returns `*XService`.
- Postgres repo constructor: `NewPosgresRepository(db)` or `NewXRepository(db)` (note: existing code spells it "Posgres" — match the surrounding file).
- Handler constructor: `NewXHandler(service, validator)` returns `*XHandler`.

## Module struct shape

The `Module` struct exposes whatever sibling modules need at wire time. Most modules only expose `Service`; `user.Module` also exposes `Repo` because `auth.NewModule` consumes it (see `internal/user/routes.go:11-14` and `internal/auth/routes.go:15`).

Only export the minimum surface needed by `cmd/web/routes.go`.
