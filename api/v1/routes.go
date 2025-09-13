package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/iragsraghu/user-service/internal/handler"
)

func Register(r *gin.Engine, h *handler.UserHandler) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/users", h.Create)
		v1.GET("/users/:id", h.GetByID)
		v1.GET("/users", h.List)
		v1.GET("/health", h.Health)
	}
}
