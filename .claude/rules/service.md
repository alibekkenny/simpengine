# Service conventions

Services hold all business logic. They are the only layer allowed to:

- Read auth from the request context.
- Coordinate multiple repositories.
- Call other modules' services.
- Translate repo errors into the domain-level sentinels handlers expect.

## Constructor

Repositories are accepted as **interface types** (defined in the same module's `repository.go`), never as concrete structs:

```go
type XService struct {
    repo XRepository    // interface from repository.go
}

func NewXService(r XRepository) *XService {
    return &XService{repo: r}
}
```

Cross-module dependencies are accepted as **services**, never repositories:

```go
func NewRomanticEventService(
    repo repository.RomaticEventRepository,
    stepRepo repository.EventStepRepository,
    optionRepo repository.EventStepOptionRepository,
    simpTargetService *simptarget.SimpTargetService,
    mediaService *media.MediaService,
    userService *user.UserService,
    notifier *notification.NotificationService,
) *RomanticEventService { ... }
```

Reference: `internal/romantic_event/service.go:30-47`.

## Auth context

Every authenticated service method must pull the user id off the context as its **first** step:

```go
userID, ok := auth.GetUserIDFromContext(ctx)
if !ok {
    return 0, model.ErrInvalidCredentials   // or `return nil, ...` / `return ...`
}
```

Reference: `internal/simp-target/service.go:21-24`. Never trust a userID passed in from the handler — always read it from context inside the service.

## Authorization check pattern

When the action requires another entity to be owned by the user, fetch via the "by id and user" repo method first; the repo returns `model.ErrNoRecord` if either the row doesn't exist or it belongs to someone else (don't leak the distinction):

```go
event, err := s.repo.FindByIDAndUserID(ctx, id, userID)
if err != nil {
    if errors.Is(err, model.ErrNoRecord) {
        return fmt.Errorf("%w: romantic event not found", model.ErrNoRecord)
    }
    return fmt.Errorf("%w: %v", model.ErrInternal, err)
}
```

Reference: `internal/romantic_event/service.go:80-89`.

## State-machine guards

When an entity has a status field, enforce transitions in the service with `ErrInvalidState`:

```go
if event.Status != rmodel.StatusDraft {
    return fmt.Errorf("%w: cannot edit event with status %s", model.ErrInvalidState, event.Status)
}
```

Reference: `internal/romantic_event/service.go:87-89`.

## Error wrapping

Always wrap with `%w` to preserve the sentinel:

```go
return fmt.Errorf("%w: %v", model.ErrInternal, err)        // unknown repo error
return fmt.Errorf("%w: simp target not found", model.ErrNoRecord)  // rephrase known sentinel
return fmt.Errorf("%w: invalid login format", model.ErrInvalidBody) // domain validation
```

Never return raw repo errors to the handler — they should always be wrapped in a sentinel. See [error-handling.md](error-handling.md).

## Cross-module calls

Call other modules through their service, not their repo. Pre-checks on related entities go through the owning service so ownership is enforced uniformly:

```go
// in romantic_event.service.go
if _, err := s.simpTargetService.GetSimpTargetByIDAndUser(ctx, simpTargetID); err != nil {
    return 0, err   // already wrapped
}
```

Reference: `internal/romantic_event/service.go:55-57`.
