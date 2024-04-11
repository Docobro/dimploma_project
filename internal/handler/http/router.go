package http

import (
	"github.com/docobro/dimploma_project/pkg/logger"
	"github.com/gin-gonic/gin"
)

type HandlerRouter interface {
	AddRoutes(r *gin.RouterGroup)
	GetVersion() string
	GetContentType() string
}

type Router struct {
	router *gin.Engine
}

func NewRouter() *Router {
	return &Router{router: gin.Default()}
}

func (r *Router) WithMetrics() *Router {
	return r
}

func (r *Router) WithHandler(h HandlerRouter, logger logger.Logger) *Router {
	r.router.Use(gin.Recovery())
	// Routers
	api := r.router.Group("/api/" + h.GetVersion())
	// Uncomment below lines if you have middleware functions
	// api.Use(middleware.AddContextMiddleware(logger))
	// r.router.Use(middleware.AccessLogMiddleware(logger))
	h.AddRoutes(api)

	return r
}
