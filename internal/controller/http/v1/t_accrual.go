package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// TESTAccrual test external service handler. Always return "OK" and "100" as response
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
	ctx.String(http.StatusOK, "100")
}
