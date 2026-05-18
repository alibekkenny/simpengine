---
description: Scaffold a new internal/<module> following the controller/service/repository pattern
argument-hint: <module-name>
---

Scaffold a new module named `$ARGUMENTS` under `internal/`.

Before writing any code, read these rule files so you follow the conventions exactly:

- @.claude/rules/module-structure.md
- @.claude/rules/handler.md
- @.claude/rules/service.md
- @.claude/rules/repository.md
- @.claude/rules/error-handling.md
- @.claude/rules/auth.md
- @.claude/rules/wiring.md

Then look at an existing simple module to copy the pattern from — `internal/simp-target/` is the cleanest reference (single entity, full CRUD, owned-by-user, no cross-module deps).

## Steps

1. Confirm with me:
   - Module name and package name (folder may be hyphenated; package must be lowercase, no hyphens).
   - Primary entity (struct + SQL columns).
   - CRUD set needed (create / read / update / delete / list).
   - Whether rows are owned by a user (almost always yes) or scoped some other way.
   - Whether any sibling module's *service* is needed (pass-through ownership checks, etc.).

2. Create the files:
   ```
   internal/<module>/
   ├── routes.go         # Module struct, NewModule(db, ...deps), RegisterRoutes(mux, cfg)
   ├── handler.go        # XHandler + endpoint methods, all with Swagger annotations
   ├── service.go        # XService + business logic, reading userID from auth context
   ├── repository.go     # XRepository interface
   ├── postgres_repo.go  # PostgresRepository implementation
   ├── model.go          # X struct
   └── dto.go            # request/response DTOs with validator tags
   ```

3. Wire it into `cmd/web/routes.go`:
   - Import the new package.
   - Construct the module **after** every dep it consumes.
   - Call `RegisterRoutes(mux, app.config)`.

4. Add a migration pair under `migrations/` using `migrate create -ext sql -dir migrations -seq create_<table>_table`. Use the next zero-padded 6-digit sequence number. Don't leave the `.down.sql` empty.

5. Regenerate Swagger: `swag init -g cmd/web/main.go -o cmd/web/docs`.

6. Verify it builds: `go build -o /tmp/main ./cmd/web`.

## What MUST be true when you're done

- All handlers go through `shared_model.WriteErrorResponse(w, err)` for every error path.
- All authenticated service methods read `userID` from `auth.GetUserIDFromContext(ctx)` *before* anything else.
- All repo functions pipe DB errors through `db.MapPQError(err)`; UPDATE/DELETE check `RowsAffected` and return `model.ErrNoRecord` on 0.
- All endpoint registrations use `auth.AuthMiddleware(http.HandlerFunc(...))` unless explicitly public.
- Every handler has a full Swagger annotation block.
- `go build ./cmd/web` succeeds.
