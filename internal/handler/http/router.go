package http

import (
	"github.com/docobro/dimploma_project/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	r.router.GET("/metrics", prometheusHandler())
	return r
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
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
