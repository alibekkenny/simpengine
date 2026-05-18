# Auth

## How a request becomes authenticated

1. `POST /user/login` (in `internal/auth/handler.go`) verifies credentials and calls `auth.GenerateJWT(...)`.
2. The JWT is set on the response in a cookie:
   - Name: `jwt`
   - Flags: `HttpOnly`, `Secure`, `SameSite=None`
3. On subsequent requests, `auth.AuthMiddleware` reads the `jwt` cookie, parses it, and stores claims in the request context under unexported keys.

## Reading auth in services

Always use the accessors — never reach into the context key directly (the key is unexported on purpose):

```go
userID, ok := auth.GetUserIDFromContext(ctx)  // int64
if !ok {
    return model.ErrInvalidCredentials
}

role, ok := auth.GetRoleFromContext(ctx)      // string
```

Reference: `internal/auth/middleware.go:66-74`.

## Wiring middleware on routes

`AuthMiddleware` wraps the `http.Handler`, not the `http.HandlerFunc`. The pattern in `RegisterRoutes`:

```go
mux.Handle("POST /simp-target",
    auth.AuthMiddleware(http.HandlerFunc(handler.CreateSimpTarget)))
```

For admin-only routes, compose with `RoleMiddleware`:

```go
mux.Handle("POST /admin/template-event/steps",
    auth.AuthMiddleware(auth.RoleMiddleware("admin")(http.HandlerFunc(handler.AddTemplateEventStep))))
```

Reference: `internal/romantic_event/routes.go:35-65`.

## Public routes

Unauthenticated endpoints (e.g. the simp target's `/public/romantic-event/{public_token}` views) are registered with `mux.Handle(...)` and `http.HandlerFunc(...)` only — no `AuthMiddleware`. The service implementations for these reach data via `FindByPublicToken`, not `FindByIDAndUserID`. See `internal/romantic_event/routes.go:66-69`.

## CORS / cookies

CORS allow-origin is `cfg.FrontEndHost` (single origin) with credentials enabled. That's a hard requirement for the `jwt` cookie to be sent cross-site. Any new origin must come through that config field.

## Don't

- Don't read or write the JWT cookie outside `internal/auth/`.
- Don't add a parallel auth scheme (header, bearer). The cookie is the source of truth.
- Don't `errors.New(...)` your own "unauthorized" — return `model.ErrInvalidCredentials` so the response maps to 401 with the right machine code.
- Don't accept `userID` as a handler/service parameter. Read it from context inside the service.
