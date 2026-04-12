# Backend Guidelines

Stack: Go · PocketBase (embedded) · Gin · internal package layout

Ports: PocketBase admin + DB on **:8090** · Gin API on **:8313**

---

## Response format

**All** endpoints must return responses through `models.Respond`. Never call `c.JSON` directly.

```go
// internal/models/response.go
models.Respond(c, http.StatusOK, "ok", data)
models.Respond(c, http.StatusBadRequest, "invalid input", nil)
models.Respond(c, http.StatusInternalServerError, "something went wrong", nil)
```

Shape: `{ "message": string, "data": any | null }`

---

## Router groups

Group routes by resource. Every group lives in `internal/router/router.go`.

```go
api := r.Group("/api")
{
    api.GET("/health", handlers.Health)

    users := api.Group("/users")
    users.Use(middleware.RequireAuth(app))
    {
        users.GET("", userHandler.List)
        users.GET("/:id", userHandler.Get)
        users.POST("", userHandler.Create)
        users.PUT("/:id", userHandler.Update)
        users.DELETE("/:id", userHandler.Delete)
    }
}
```

---

## Package structure

```
cmd/api/main.go          entry point — wires everything together
internal/config/         env config loaded once at startup
internal/handlers/       one file per resource (users.go, items.go, …)
internal/middleware/     auth.go, cors.go, role checks
internal/models/         response.go + domain struct types
internal/services/       business logic — handlers call services, not Dao directly
internal/router/         route registration only, no logic
pb_migrations/           PocketBase schema migrations
```

---

## Middleware

Common checks belong in `internal/middleware/`, not in handlers.

```go
// Protect a group
protected := api.Group("/")
protected.Use(middleware.RequireAuth(app))

// Role check (see skills/backend/auth-middleware for RequireRole)
admin := api.Group("/admin")
admin.Use(middleware.RequireAuth(app), middleware.RequireRole("admin"))
```

---

## Business logic

Handlers should stay thin — they parse input, call a service, respond. All logic goes in `internal/services/`.

```go
// Handler
func (h *ItemHandler) Create(c *gin.Context) {
    var body CreateItemInput
    if err := c.ShouldBindJSON(&body); err != nil {
        models.Respond(c, http.StatusBadRequest, "invalid body", nil)
        return
    }
    record, err := h.service.Create(body)
    if err != nil {
        models.Respond(c, http.StatusInternalServerError, "failed to create", nil)
        return
    }
    models.Respond(c, http.StatusCreated, "created", record)
}
```

---

## PocketBase — schema via migrations (no admin UI needed)

**Never design collections through the admin UI.** Define all schema as Go migration files so everything is reproducible and version-controlled.

### Creating a migration

```bash
go run cmd/api/main.go migrate create add_posts_collection
# generates pb_migrations/<timestamp>_add_posts_collection.go
```

Fill in the generated file (see `pb_migrations/1_initial.go` for a full example):

```go
func init() {
    m.Register(func(db dbx.Builder) error {
        dao := daos.New(db)
        collection := &models.Collection{
            Name: "posts",
            Type: models.CollectionTypeBase,
            Schema: schema.NewSchema(
                &schema.SchemaField{Name: "title", Type: schema.FieldTypeText, Required: true},
            ),
            ListRule:   strPtr("@request.auth.id != ''"),
            CreateRule: strPtr("@request.auth.id != ''"),
        }
        return dao.SaveCollection(collection)
    }, func(db dbx.Builder) error {
        // rollback
        dao := daos.New(db)
        col, _ := dao.FindCollectionByNameOrId("posts")
        if col != nil { return dao.DeleteCollection(col) }
        return nil
    })
}
```

Migrations run automatically on `./main serve`. **Never edit an existing migration** — always create a new one.

### Admin account

Auto-created from env vars on first run (see `.env.example`). Set `PB_ADMIN_EMAIL` + `PB_ADMIN_PASSWORD` and the server bootstraps it. The admin UI at `:8090/_/` is available but only needed for manual inspection — all schema changes go through migrations.

### What PocketBase vs Gin handles

- **PocketBase**: users auth, file storage, real-time subscriptions, collection rules, admin UI
- **Gin**: custom business logic, AI integrations, aggregations, anything that needs Go code

---

## Adding a new executable

Add the entry point at `cmd/<name>/main.go`, then add to `Dockerfile`:

```dockerfile
RUN CGO_ENABLED=0 GOOS=linux go build -o <name> cmd/<name>/main.go

FROM alpine:3.20.1 AS <name>
WORKDIR /app
COPY --from=build /app/<name> /app/<name>
CMD ["./<name>"]
```

And add a service in `docker-compose.yml` with `target: <name>`.

---

## Environment variables

All config is loaded from env via `internal/config/config.go`. Never use `os.Getenv` directly in handlers or services — always use the `*Config` struct passed at startup.
