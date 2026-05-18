# simpengine rules

Detailed conventions for working in this Go backend. Load the relevant rule before writing or reviewing code in that layer.

- [module-structure.md](module-structure.md) — file layout inside `internal/<module>/`
- [handler.md](handler.md) — HTTP handlers + Swagger annotations
- [service.md](service.md) — business logic, auth context, cross-module deps
- [repository.md](repository.md) — Postgres access, error mapping, RowsAffected
- [error-handling.md](error-handling.md) — sentinel errors, HTTP status mapping, error wrapping
- [auth.md](auth.md) — JWT cookie, middleware, user/role context
- [wiring.md](wiring.md) — module construction and route registration in `cmd/web/routes.go`

See top-level `CLAUDE.md` for project overview and commands.
