package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/account"
	"net/http"
)

type decreaseBalanceRequest struct {
	UserId account.UserID `json:"user_id"`
	Sum    int            `json:"sum"`
}

type registerAccountRequest struct {
	Balance int `json:"balance" binding:"required"`
}

// RegisterAccount register a new account with custom balance.
func (h *Handler) RegisterAccount(ctx *gin.Context) {

	request := new(registerAccountRequest)
	if err := ctx.ShouldBindJSON(request); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newAccount := account.Account{
		ID:      account.UserID(uuid.NewString()),
		Balance: request.Balance,
	}
	if err := newAccount.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperr.AccountInvalidBalance})
		return
	}

	if err := h.services.Account.CreateAccount(ctx, newAccount); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok", "user_id": newAccount.ID})
}

func (h *Handler) DecreaseBalance(ctx *gin.Context) {

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
