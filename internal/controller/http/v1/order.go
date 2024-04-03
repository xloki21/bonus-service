package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xloki21/bonus-service/internal/apperr"
	t "github.com/xloki21/bonus-service/internal/entity/order"
	"net/http"
)

func (h *Handler) RegisterOrder(ctx *gin.Context) {
	var order = new(t.Order)

	if err := ctx.ShouldBindJSON(order); err != nil {

	}

	if err := order.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": apperr.OrderValidationFailed})
		return
	}

	if err := h.services.Order.Register(ctx, order); err != nil {
		if errors.Is(err, apperr.OrderAlreadyRegistered) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Order registered"})
}
