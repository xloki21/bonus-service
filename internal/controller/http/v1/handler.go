package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xloki21/bonus-service/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) ApiV1(mode string) *gin.Engine {
	gin.SetMode(mode)
	router := gin.New()

	//TEST ACCRUAL SERVER
	router.GET("/info", h.TESTAccrual)

	api := router.Group("/api/v1")
	{
		api.GET("/ping", h.Ping)
		api.POST("/register", h.RegisterOrder)
		api.POST("/decrease", h.Decrease)
	}

	return router
}
