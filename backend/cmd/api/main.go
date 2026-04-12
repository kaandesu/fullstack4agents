package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	_ "fullstack4agents/pb_migrations" // register all migrations
	"fullstack4agents/internal/config"
	"fullstack4agents/internal/router"
)

func main() {
	cfg := config.Load()
	app := pocketbase.New()

	// Register the migrate command so `./main migrate` works
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// Auto-create migration files when running `./main migrate create <name>`
		Automigrate: true,
	})

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// Bootstrap admin from env on first run (no-op if an admin already exists)
		bootstrapAdmin(app)

		// Start Gin API server
		r := router.New(app)
		go func() {
			if err := r.Run(":" + cfg.Port); err != nil {
				log.Fatal("Gin server error:", err)
			}
		}()
		return nil
	})

	// PocketBase blocks here — runs migrations automatically, serves admin UI on :8090
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// bootstrapAdmin creates a superadmin from env vars if that email doesn't exist yet.
// Set PB_ADMIN_EMAIL and PB_ADMIN_PASSWORD in your .env on first deploy.
func bootstrapAdmin(app *pocketbase.PocketBase) {
	email := os.Getenv("PB_ADMIN_EMAIL")
	password := os.Getenv("PB_ADMIN_PASSWORD")
	if email == "" || password == "" {
		return
	}
	// No-op if an admin with this email already exists
	if _, err := app.Dao().FindAdminByEmail(email); err == nil {
		return
	}
	admin := &models.Admin{}
	admin.Email = email
	if err := admin.SetPassword(password); err != nil {
		log.Println("Admin bootstrap: failed to set password:", err)
		return
	}
	if err := app.Dao().SaveAdmin(admin); err != nil {
		log.Println("Admin bootstrap: failed to save admin:", err)
		return
	}
	log.Println("Admin bootstrap: created admin account for", email)
}
