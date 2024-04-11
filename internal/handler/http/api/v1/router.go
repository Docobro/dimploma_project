package v1

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetVersion() string {
	return "v1"
}

func (h *Handler) GetContentType() string {
	return ""
}

func (h *Handler) AddRoutes(r *gin.RouterGroup) {
	r.GET("/ping", h.Ping)
}
