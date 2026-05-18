# Error handling

All errors are funneled through a small set of sentinels defined in `internal/shared/model/error.go`. The HTTP layer never decides on a status code — it asks the sentinel via `WriteErrorResponse`.

## Sentinels and their HTTP mapping

| Sentinel                | HTTP | Code                  | Use it when…                                                       |
| ----------------------- | ---- | --------------------- | ------------------------------------------------------------------ |
| `ErrInvalidBody`        | 400  | `INVALID_BODY`        | JSON decode failed, or body field has invalid shape/format.        |
| `ErrInvalidParams`      | 400  | `INVALID_PARAMETERS`  | Path or query param couldn't be parsed (e.g. non-numeric id).      |
| `ErrValidation`         | 400  | `VALIDATION_ERROR`    | `validator.Struct(...)` returned an error.                         |
| `ErrInvalidCredentials` | 401  | `INVALID_CREDENTIALS` | No user id in context, bad password, bad/missing JWT.              |
| `ErrForbidden`          | 403  | `FORBIDDEN`           | Authenticated but not allowed (use sparingly; usually 404 instead).|
| `ErrNoRecord`           | 404  | `NOT_FOUND`           | Row missing, or row exists but isn't owned by the caller.          |
| `ErrUniqueViolation`    | 409  | `UNIQUE_VIOLATION`    | Duplicate insert (unique index hit).                               |
| `ErrInvalidState`       | 409  | `INVALID_STATE`       | Entity is in a state that disallows this transition (e.g. published event being edited). |
| `ErrInternal`           | 500  | `INTERNAL_ERROR`      | Anything else. Underlying error is logged but not sent to client.  |

Reference: `internal/shared/model/error.go`.

## Wrapping rules

Use `%w` so `errors.Is(err, model.ErrX)` still works after wrapping. Always lead with the sentinel.

```go
return fmt.Errorf("%w: %v", model.ErrInternal, err)             // unknown failure
return fmt.Errorf("%w: user not found", model.ErrNoRecord)      // rephrased domain error
return fmt.Errorf("%w:\n%v", model.ErrValidation, err)          // attaches validator details
```

Use `\n` (not `: `) when attaching multi-line `validator` output — matches the rest of the codebase (`internal/user/handler.go:36,41`).

Never return a raw external error (`*pq.Error`, `sql.ErrNoRows`, `bcrypt.*`, etc.) to the handler. Translate at the lowest layer that knows how:

- `sql.ErrNoRows` → `model.ErrNoRecord` (in the repo).
- `*pq.Error` → through `db.MapPQError` (in the repo).
- Anything still bare reaching the service → wrap as `ErrInternal`.

## Handler error path

Handlers do **not** branch on error kind. They write whatever they get and let `WriteErrorResponse` map it:

```go
if err := h.service.Do(ctx, ...); err != nil {
    shared_model.WriteErrorResponse(w, err)
    return
}
```

The only errors created in handlers are body/param parsing errors, which are wrapped at the site of failure:

```go
shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
```

## Don't add a new sentinel without thinking

Prefer reusing the existing ones. A new sentinel needs:

1. A `var Err... = errors.New(...)` in `internal/shared/model/error.go`.
2. A new `case` in `ErrorStatus` returning the HTTP status and machine code.
3. Updated Swagger `@Failure` lines in every handler that can produce it.
