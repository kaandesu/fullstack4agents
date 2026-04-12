# Skill: CRUD Handler

A complete Create / Read / Update / Delete handler for a PocketBase collection.

Copy to `backend/internal/handlers/<resource>.go` and adapt the collection name and struct.

---

## Handler

```go
// backend/internal/handlers/posts.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pocketbase/pocketbase"
	"fullstack4agents/internal/models"
)

type PostHandler struct {
	app *pocketbase.PocketBase
}

func NewPostHandler(app *pocketbase.PocketBase) *PostHandler {
	return &PostHandler{app: app}
}

func (h *PostHandler) List(c *gin.Context) {
	records, err := h.app.Dao().FindRecordsByFilter("posts", "1=1", "-created", 50, 0)
	if err != nil {
		models.Respond(c, http.StatusInternalServerError, "failed to fetch posts", nil)
		return
	}
	models.Respond(c, http.StatusOK, "ok", records)
}

func (h *PostHandler) Get(c *gin.Context) {
	record, err := h.app.Dao().FindRecordById("posts", c.Param("id"))
	if err != nil {
		models.Respond(c, http.StatusNotFound, "post not found", nil)
		return
	}
	models.Respond(c, http.StatusOK, "ok", record)
}

func (h *PostHandler) Create(c *gin.Context) {
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		models.Respond(c, http.StatusBadRequest, "invalid body", nil)
		return
	}
	collection, err := h.app.Dao().FindCollectionByNameOrId("posts")
	if err != nil {
		models.Respond(c, http.StatusInternalServerError, "collection not found", nil)
		return
	}
	record := models.NewRecord(collection)
	for k, v := range body {
		record.Set(k, v)
	}
	if err := h.app.Dao().SaveRecord(record); err != nil {
		models.Respond(c, http.StatusInternalServerError, "failed to create post", nil)
		return
	}
	models.Respond(c, http.StatusCreated, "created", record)
}

func (h *PostHandler) Update(c *gin.Context) {
	record, err := h.app.Dao().FindRecordById("posts", c.Param("id"))
	if err != nil {
		models.Respond(c, http.StatusNotFound, "post not found", nil)
		return
	}
	var body map[string]any
	if err := c.ShouldBindJSON(&body); err != nil {
		models.Respond(c, http.StatusBadRequest, "invalid body", nil)
		return
	}
	for k, v := range body {
		record.Set(k, v)
	}
	if err := h.app.Dao().SaveRecord(record); err != nil {
		models.Respond(c, http.StatusInternalServerError, "failed to update post", nil)
		return
	}
	models.Respond(c, http.StatusOK, "updated", record)
}

func (h *PostHandler) Delete(c *gin.Context) {
	record, err := h.app.Dao().FindRecordById("posts", c.Param("id"))
	if err != nil {
		models.Respond(c, http.StatusNotFound, "post not found", nil)
		return
	}
	if err := h.app.Dao().DeleteRecord(record); err != nil {
		models.Respond(c, http.StatusInternalServerError, "failed to delete post", nil)
		return
	}
	models.Respond(c, http.StatusOK, "deleted", nil)
}
```

---

## Registering the routes

In `backend/internal/router/router.go`:

```go
posts := api.Group("/posts")
posts.Use(middleware.RequireAuth(app))
{
	h := handlers.NewPostHandler(app)
	posts.GET("", h.List)
	posts.GET("/:id", h.Get)
	posts.POST("", h.Create)
	posts.PUT("/:id", h.Update)
	posts.DELETE("/:id", h.Delete)
}
```
