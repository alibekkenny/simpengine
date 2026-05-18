# Wiring (`cmd/web/routes.go`)

## Construction order

`cmd/web/routes.go::routes()` builds every module up-front, then asks each to register its routes. Order matters because later modules consume earlier ones' services.

Current order (extend it, don't reorder existing entries):

```go
userModule          := user.NewModule(app.config.DB)
authModule          := auth.NewModule(userModule.Repo)
simpTargetModule    := simptarget.NewModule(app.config.DB)
mediaModule         := media.NewModule(app.config.DB, app.config.MinioClient, app.config.MinioBucketName)
notificationModule  := notification.NewModule(app.config.TelegramBotToken)
romanticEventModule := romanticevent.NewModule(
    app.config.DB,
    simpTargetModule.Service,
    mediaModule.Service,
    userModule.Service,
    notificationModule.Service,
)
```

Reference: `cmd/web/routes.go:20-25`.

## Cross-module dependency rule

Pass **services**, not repositories. The one exception is `auth.NewModule(userModule.Repo)` — auth only needs the user repository to look users up by login; it does not need higher-level user business logic.

If you find yourself wanting another module's repo, instead add a thin service method on the owning module and call that.

## Register routes

After constructing all modules:

```go
userModule.RegisterRoutes(mux, app.config)
authModule.RegisterRoutes(mux, app.config)
simpTargetModule.RegisterRoutes(mux, app.config)
romanticEventModule.RegisterRoutes(mux, app.config)
mediaModule.RegisterRoutes(mux, app.config)
```

Notifications and admin invite have no HTTP routes today, so they don't call `RegisterRoutes`.

## Middleware chain

Global middleware is composed with `justinas/alice`:

```go
standardChain := alice.New(app.recoverPanic, app.logRequest, app.secureHeaders)
return standardChain.Then(mux)
```

Per-route middleware (auth, role) is applied inside the module's `RegisterRoutes` — see [auth.md](auth.md).

## Adding a new module — checklist

1. Create `internal/<module>/` following [module-structure.md](module-structure.md).
2. Add `NewModule(deps...)` and `RegisterRoutes(mux, cfg)` in `routes.go`.
3. In `cmd/web/routes.go`:
   - Add the import.
   - Construct the module **after** every dep it requires.
   - Call `RegisterRoutes` if it exposes endpoints.
4. If the module reads new config, add the field to `cmd/config/config.go` and load it in `cmd/web/main.go`. Add the env var with a dev default in `docker-compose.yml`.
5. If the module needs DB tables, write a `migrations/000NNN_*.up.sql` + `.down.sql` pair — see [repository.md](repository.md).
6. Regenerate Swagger: `swag init -g cmd/web/main.go -o cmd/web/docs`.

## Config (`cmd/config/config.go`)

`config.Config` is the single struct passed into every `RegisterRoutes(mux, cfg)`. Env vars consumed in `cmd/web/main.go`:
`ADDR`, `DSN`, `JWT_SECRET`, `MINIO_ENDPOINT`, `MINIO_BUCKET`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `FRONTEND_HOST`, `TELEGRAM_BOT_TOKEN`.

When introducing a new env var:

- Read it in `main.go` (with a `log.Fatal` if it's required).
- Add a field on `config.Config`.
- Add a dev value in `docker-compose.yml` under the `web` service's `environment:` block.
