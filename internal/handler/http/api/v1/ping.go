package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type pingResponse struct {
	Message string `json:"message"`
}

func (h Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, pingResponse{Message: "ok"})
}
