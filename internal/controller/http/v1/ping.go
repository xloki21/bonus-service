package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) Ping(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Pong!"})
}
