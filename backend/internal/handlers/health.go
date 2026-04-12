package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fullstack4agents/internal/models"
)

func Health(c *gin.Context) {
	models.Respond(c, http.StatusOK, "ok", gin.H{"status": "healthy"})
}
