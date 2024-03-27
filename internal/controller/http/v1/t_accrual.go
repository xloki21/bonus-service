package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

// TESTAccrual test external service handler. Always return "OK" and random int from [0, 10000) as response
func (h *Handler) TESTAccrual(ctx *gin.Context) {
	//ts, err := strconv.Atoi(ctx.Query("timestamp"))
	//if err != nil {
	//	return
	//}
	//_ = transaction.Transaction{
	//	UserID:       account.UserID(ctx.Query("user_id")),
	//	GoodID:       order.GoodID(ctx.Query("good_id")),
	//	Status:       transaction.UNPROCESSED,
	//	Timestamp:    int64(ts),
	//	RegisteredAt: time.Now().Unix(),
	//}
	ctx.String(http.StatusOK, fmt.Sprintf("%d", rand.Intn(10000)))
}
