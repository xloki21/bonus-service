package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

// TESTAccrual test external service handler. If "OK" returns random int from [0, 10000) as response.
func (h *Handler) TESTAccrual(ctx *gin.Context) {

	// simulate accrual not found yet
	condition := rand.Intn(100)
	if condition >= 80 && condition < 90 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}

	// simulate slow request
	if condition >= 90 {
		time.Sleep(time.Second * time.Duration(rand.Intn(10)))
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("%d", rand.Intn(10000)))
}
