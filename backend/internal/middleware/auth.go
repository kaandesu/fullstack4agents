package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pocketbase/pocketbase"
	"fullstack4agents/internal/models"
)

// RequireAuth validates a PocketBase Bearer token and sets "authRecord" in context.
func RequireAuth(app *pocketbase.PocketBase) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			models.Respond(c, http.StatusUnauthorized, "missing token", nil)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		record, err := app.Dao().FindAuthRecordByToken(token, app.Settings().RecordAuthToken.Secret)
		if err != nil {
			models.Respond(c, http.StatusUnauthorized, "invalid token", nil)
			c.Abort()
			return
		}
		c.Set("authRecord", record)
		c.Next()
	}
}
