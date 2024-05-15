package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"net/http"
)

type decreaseBalanceRequest struct {
	UserId string `json:"user_id" binding:"required"`
	Sum    uint   `json:"sum" binding:"required"`
}

type registerAccountRequest struct {
	Balance uint `json:"balance" binding:"required"`
}

// RegisterAccount register a new account with custom balance.
func (h *Handler) RegisterAccount(ctx *gin.Context) {

	request := &registerAccountRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("bad request: %w", err).Error()})
		return
	}

	newAccount := account.Account{
		ID:      uuid.NewString(),
		Balance: request.Balance,
	}

	if err := h.services.Account.CreateAccount(ctx, newAccount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "user_id": newAccount.ID})
}

func (h *Handler) DecreaseBalance(ctx *gin.Context) {

	request := new(decreaseBalanceRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("bad request: %w", err).Error()})
		return
	}

	err := h.services.Account.Debit(ctx, request.UserId, request.Sum)
	if err != nil {
		if errors.Is(err, apperr.InsufficientBalance) {
			ctx.JSON(http.StatusPaymentRequired, gin.H{"error": apperr.InsufficientBalance.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}
