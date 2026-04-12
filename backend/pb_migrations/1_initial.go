package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

// Example migration — replace or extend with your actual collections.
// Agents: create new migrations with `./main migrate create <name>` which
// auto-generates a numbered file in pb_migrations/. Never edit old migrations.
func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		// Example: "posts" collection
		posts := &models.Collection{
			Name: "posts",
			Type: models.CollectionTypeBase,
			Schema: schema.NewSchema(
				&schema.SchemaField{
					Name:     "title",
					Type:     schema.FieldTypeText,
					Required: true,
				},
				&schema.SchemaField{
					Name: "content",
					Type: schema.FieldTypeEditor,
				},
				&schema.SchemaField{
					Name: "owner",
					Type: schema.FieldTypeRelation,
					Options: &schema.RelationOptions{
						CollectionId: "_pb_users_auth_",
						MaxSelect:    func() *int { n := 1; return &n }(),
					},
				},
			),
			// Only the record owner can write; anyone authenticated can read
			ListRule:   strPtr("owner = @request.auth.id"),
			ViewRule:   strPtr("owner = @request.auth.id"),
			CreateRule: strPtr("@request.auth.id != ''"),
			UpdateRule: strPtr("owner = @request.auth.id"),
			DeleteRule: strPtr("owner = @request.auth.id"),
		}

		return dao.SaveCollection(posts)
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
