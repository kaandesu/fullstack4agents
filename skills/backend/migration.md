# Skill: PocketBase Migration

How to create and write database schema migrations. All schema changes must go through migration files — never through the admin UI.

---

## Creating a migration

```bash
make migrate name=add_posts_collection
# or directly:
cd backend && go run cmd/api/main.go migrate create add_posts_collection
```

This generates `backend/pb_migrations/<timestamp>_add_posts_collection.go`. Fill in the `up` and `down` functions.

Migrations run automatically every time the server starts. **Never edit an existing migration** — always create a new one.

---

## Migration template

```go
package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection := &models.Collection{
			Name: "posts",
			Type: models.CollectionTypeBase, // or CollectionTypeAuth for user collections
			Schema: schema.NewSchema(
				// --- Common field types ---
				&schema.SchemaField{
					Name:     "title",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name: "content",
					Type: schema.FieldTypeEditor, // rich text
				},
				&schema.SchemaField{
					Name: "published",
					Type: schema.FieldTypeBool,
				},
				&schema.SchemaField{
					Name: "views",
					Type: schema.FieldTypeNumber,
				},
				&schema.SchemaField{
					Name: "tags",
					Type: schema.FieldTypeJson, // array/object stored as JSON
				},
				&schema.SchemaField{
					Name: "cover",
					Type: schema.FieldTypeFile,
					Options: &schema.FileOptions{
						MaxSelect: 1,
						MaxSize:   5242880, // 5MB
					},
				},
				&schema.SchemaField{
					Name: "owner",
					Type: schema.FieldTypeRelation,
					Options: &schema.RelationOptions{
						CollectionId: "_pb_users_auth_",
						MaxSelect:    intPtr(1),
					},
				},
			),
			// Access rules — empty string = public, nil = admin only
			ListRule:   strPtr("owner = @request.auth.id"),
			ViewRule:   strPtr("owner = @request.auth.id"),
			CreateRule: strPtr("@request.auth.id != ''"), // any authenticated user
			UpdateRule: strPtr("owner = @request.auth.id"),
			DeleteRule: strPtr("owner = @request.auth.id"),
		}

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		// Down — rollback
		dao := daos.New(db)
		col, _ := dao.FindCollectionByNameOrId("posts")
		if col != nil {
			return dao.DeleteCollection(col)
		}
		return nil
	})
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
```

---

## Adding a field to an existing collection

Create a new migration — do not edit the previous one:

```go
func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("posts")
		if err != nil {
			return err
		}

		collection.Schema.AddField(&schema.SchemaField{
			Name: "slug",
			Type: schema.FieldTypeText,
		})

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)
		collection, err := dao.FindCollectionByNameOrId("posts")
		if err != nil {
			return err
		}
		collection.Schema.RemoveField(collection.Schema.GetFieldByName("slug").Id)
		return dao.SaveCollection(collection)
	})
}
```

---

## Access rule reference

| Rule value | Meaning |
|---|---|
| `nil` | Admin-only access |
| `""` (empty string) | Public — anyone can access |
| `"@request.auth.id != ''"` | Any authenticated user |
| `"owner = @request.auth.id"` | Only the record owner |
| `"role = 'admin'"` | Users with role field = admin |
