package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xloki21/bonus-service/internal/apperr"
	"github.com/xloki21/bonus-service/internal/entity/order"
	"net/http"
)

func (h *Handler) RegisterOrder(ctx *gin.Context) {
	var newOrder order.Order

	if err := ctx.ShouldBindJSON(&newOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("bad request: %w", err).Error()})
		return
	}

	if err := newOrder.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("request validation error: %w", err).Error()})
		return
	}

	if err := h.services.Order.Register(ctx, newOrder); err != nil {
		if errors.Is(err, apperr.OrderAlreadyRegistered) {
			ctx.JSON(http.StatusConflict, gin.H{"error": apperr.OrderAlreadyRegistered.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Order registered"})
}
