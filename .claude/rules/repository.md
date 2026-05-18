# Repository conventions

Repositories own SQL only. No business logic, no auth, no error rephrasing — those belong in the service.

## Interface goes in `repository.go`

```go
package simptarget

import "context"

type SimpTargetRepository interface {
    CreateSimpTarget(ctx context.Context, name, description string, userID int64) (int64, error)
    UpdateSimpTarget(ctx context.Context, id int64, name, description string, userID int64) error
    DeleteSimpTarget(ctx context.Context, id int64, userID int64) error
    FindAllByUserID(ctx context.Context, userID int64) ([]*SimpTarget, error)
    FindByIDAndUserID(ctx context.Context, id, userID int64) (*SimpTarget, error)
    FindByID(ctx context.Context, id int64) (*SimpTarget, error)
}
```

Reference: `internal/simp-target/repository.go`.

## Postgres impl lives in `postgres_repo.go` (or `repository/postgres/<entity>.go`)

```go
type PostgresRepository struct {
    db *sql.DB
}

func NewPosgresRepository(db *sql.DB) *PostgresRepository {
    return &PostgresRepository{db: db}
}
```

For multi-entity modules, each entity gets its own file under `repository/postgres/`, each defining its own `XRepository` struct and `NewXRepository(db)`. See `internal/romantic_event/repository/postgres/`.

## SQL patterns

Use `database/sql` directly with `context.Context`:

```go
// INSERT … RETURNING id
err := r.db.QueryRowContext(ctx, stmt, args...).Scan(&id)

// SELECT single row
err := r.db.QueryRowContext(ctx, stmt, args...).Scan(&dest...)

// SELECT multiple rows
rows, err := r.db.QueryContext(ctx, stmt, args...)
defer rows.Close()
for rows.Next() { ... }

// UPDATE / DELETE
result, err := r.db.ExecContext(ctx, stmt, args...)
```

## Mandatory error mapping

Every error returned from a `database/sql` call MUST flow through `db.MapPQError(err)` before reaching the caller. This converts `*pq.Error` codes into our domain sentinels (unique violation → `ErrUniqueViolation`, FK violation → `ErrNoRecord`).

```go
err := r.db.QueryRowContext(ctx, stmt, args...).Scan(...)
if err != nil {
    return 0, db.MapPQError(err)
}
```

Reference: `internal/shared/db/helpers.go`. Never return a raw `*pq.Error`.

## "Not found" for single-row reads

`sql.ErrNoRows` becomes `model.ErrNoRecord`:

```go
err := r.db.QueryRowContext(ctx, stmt, id).Scan(&dest...)
if err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return nil, model.ErrNoRecord
    }
    return nil, db.MapPQError(err)
}
```

## "Not found" for UPDATE / DELETE

`ExecContext` does not return `sql.ErrNoRows`. Always check `RowsAffected` and translate 0 to `model.ErrNoRecord`:

```go
result, err := r.db.ExecContext(ctx, stmt, args...)
if err != nil {
    return db.MapPQError(err)
}
rowsAffected, err := result.RowsAffected()
if err != nil {
    return db.MapPQError(err)
}
if rowsAffected == 0 {
    return model.ErrNoRecord
}
return nil
```

Reference: `internal/simp-target/postgres_repo.go:33-50`.

## Ownership in queries

Most "by id" lookups also constrain by `user_id` so unauthorized rows look identical to missing rows:

```sql
SELECT id, name, description FROM simp_targets WHERE id = $1 AND user_id = $2
```

Pair these with a `FindByIDAndUserID` repo method and call it from the service after pulling `userID` from context. See [service.md](service.md).

## What goes in `model.go`

Plain structs that mirror DB rows. JSON tags are fine here (the model is also serialized to clients in many endpoints). Sensitive fields use `json:"-"` (see `internal/user/model.go:11` for the password hash).

## Migrations

- Files live in `migrations/`, named `000NNN_<snake_case>.up.sql` + `.down.sql`.
- Sequence number is zero-padded to 6 digits and monotonic.
- Create with `migrate create -ext sql -dir migrations -seq <name>`.
- Migrations run on container start via `entrypoint.sh`; no manual step needed for the dockerized DB.
- Always write a real `.down.sql` — don't leave it empty.
