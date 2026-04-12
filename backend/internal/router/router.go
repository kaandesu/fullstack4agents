package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pocketbase/pocketbase"
	"fullstack4agents/internal/handlers"
	"fullstack4agents/internal/middleware"
)

func New(app *pocketbase.PocketBase) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)

		// Protected routes — add handlers here using middleware.RequireAuth(app)
		protected := api.Group("/")
		protected.Use(middleware.RequireAuth(app))
		{
			// Example: protected.GET("/me", handlers.Me)
		}
	}

	return r
}
