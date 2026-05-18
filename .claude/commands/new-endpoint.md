---
description: Add a new HTTP endpoint to an existing module (handler + service + repo method)
argument-hint: <module> <METHOD> <path>
---

Add an endpoint to an existing module. Argument format: `<module> <METHOD> <path>`, e.g. `simp-target POST /simp-target/{id}/archive`.

Before writing code, read these rule files:

- @.claude/rules/handler.md
- @.claude/rules/service.md
- @.claude/rules/repository.md
- @.claude/rules/error-handling.md
- @.claude/rules/auth.md

## Steps

1. Open the target module under `internal/<module>/`. Skim its existing `handler.go`, `service.go`, `repository.go` so the new method matches the style of its neighbors (parameter ordering, naming, error wrapping).

2. If the change needs new DB access:
   - Add the method to the repo interface in `repository.go`.
   - Implement it in `postgres_repo.go` (or `repository/postgres/<entity>.go`).
   - Follow the patterns in @.claude/rules/repository.md (MapPQError, RowsAffected, sql.ErrNoRows → ErrNoRecord).

3. Add the service method:
   - First line should be `userID, ok := auth.GetUserIDFromContext(ctx)` for any authenticated path.
   - Wrap repo errors with `model.ErrNoRecord` / `model.ErrInternal` etc. via `fmt.Errorf("%w: ...", ...)`.
   - For status transitions, guard with `model.ErrInvalidState`.

4. Add the handler:
   - Parse path params with `r.PathValue(...)` → `ErrInvalidParams` on bad input.
   - Decode body → `ErrInvalidBody`, then `validator.Struct(...)` → `ErrValidation`.
   - Call the service; pipe any error through `shared_model.WriteErrorResponse(w, err)`.
   - Set `Content-Type` before `WriteHeader`. Use 201 + `Location` for creates, 204 for `PUT`/`DELETE` without a body, 200 for reads.
   - Write the full Swagger annotation block.

5. Register the route in `routes.go` with the right middleware (`auth.AuthMiddleware`, plus `auth.RoleMiddleware("admin")` if applicable).

6. Regenerate Swagger: `swag init -g cmd/web/main.go -o cmd/web/docs`.

7. Verify: `go build -o /tmp/main ./cmd/web`.
