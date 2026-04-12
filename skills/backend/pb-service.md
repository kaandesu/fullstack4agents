# Skill: PocketBase Service

A thin service wrapper around a PocketBase collection. Handlers stay thin by calling the service instead of the Dao directly.

Copy to `backend/internal/services/<resource>_service.go`.

---

## Service

```go
// backend/internal/services/posts_service.go
package services

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

type PostsService struct {
	app *pocketbase.PocketBase
}

func NewPostsService(app *pocketbase.PocketBase) *PostsService {
	return &PostsService{app: app}
}

func (s *PostsService) FindAll() ([]*models.Record, error) {
	return s.app.Dao().FindRecordsByFilter("posts", "1=1", "-created", 100, 0)
}

func (s *PostsService) FindByID(id string) (*models.Record, error) {
	return s.app.Dao().FindRecordById("posts", id)
}

func (s *PostsService) FindByOwner(ownerID string) ([]*models.Record, error) {
	return s.app.Dao().FindRecordsByFilter("posts", "owner = {:owner}", "-created", 100, 0,
		map[string]any{"owner": ownerID})
}

func (s *PostsService) Create(data map[string]any) (*models.Record, error) {
	col, err := s.app.Dao().FindCollectionByNameOrId("posts")
	if err != nil {
		return nil, err
	}
	record := models.NewRecord(col)
	for k, v := range data {
		record.Set(k, v)
	}
	return record, s.app.Dao().SaveRecord(record)
}

func (s *PostsService) Update(id string, data map[string]any) (*models.Record, error) {
	record, err := s.app.Dao().FindRecordById("posts", id)
	if err != nil {
		return nil, err
	}
	for k, v := range data {
		record.Set(k, v)
	}
	return record, s.app.Dao().SaveRecord(record)
}

func (s *PostsService) Delete(id string) error {
	record, err := s.app.Dao().FindRecordById("posts", id)
	if err != nil {
		return err
	}
	return s.app.Dao().DeleteRecord(record)
}
```

---

## Wiring the service into a handler

```go
// In the handler struct, hold the service instead of the app directly
type PostHandler struct {
	service *services.PostsService
}

func NewPostHandler(app *pocketbase.PocketBase) *PostHandler {
	return &PostHandler{service: services.NewPostsService(app)}
}

func (h *PostHandler) List(c *gin.Context) {
	records, err := h.service.FindAll()
	if err != nil {
		models.Respond(c, http.StatusInternalServerError, "failed to fetch posts", nil)
		return
	}
	models.Respond(c, http.StatusOK, "ok", records)
}
```
