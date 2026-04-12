# Skill: Auth & Role Middleware

The base `RequireAuth` middleware is already implemented in `backend/internal/middleware/auth.go`. This skill covers how to use it and how to add role-based checks on top.

---

## RequireAuth

Validates the PocketBase Bearer token and injects the `authRecord` into the Gin context. All protected route groups should use this.

```go
// In router.go
protected := api.Group("/")
protected.Use(middleware.RequireAuth(app))
{
    protected.GET("/me", handlers.Me)
}
```

Accessing the authenticated user inside a handler:

```go
func Me(c *gin.Context) {
    record := c.MustGet("authRecord").(*models.Record)
    models.Respond(c, http.StatusOK, "ok", map[string]any{
        "id":    record.Id,
        "email": record.GetString("email"),
    })
}
```

---

## RequireRole

Add this to `backend/internal/middleware/auth.go` when you need role-based access. Assumes a `role` field on your `users` collection.

```go
// RequireRole checks that the authenticated user has the given role value.
// Must be used after RequireAuth.
// Usage: admin.Use(middleware.RequireAuth(app), middleware.RequireRole("admin"))
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		record, exists := c.Get("authRecord")
		if !exists {
			models.Respond(c, http.StatusUnauthorized, "unauthenticated", nil)
			c.Abort()
			return
		}
		pbRecord := record.(*models.Record)
		if pbRecord.GetString("role") != role {
			models.Respond(c, http.StatusForbidden, "insufficient permissions", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
```

Example usage:

```go
admin := api.Group("/admin")
admin.Use(middleware.RequireAuth(app), middleware.RequireRole("admin"))
{
    admin.GET("/users", handlers.ListUsers)
}
```
