package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"net/http"
)

type decreaseBalanceRequest struct {
	UserId account.UserID `json:"user_id"`
	Sum    int            `json:"sum"`
}

func (h *Handler) Decrease(ctx *gin.Context) {

	request := new(decreaseBalanceRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Account.Debit(ctx, request.UserId, request.Sum)
	if err != nil {
		if errors.Is(err, apperr.InsufficientBalance) {
			ctx.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
