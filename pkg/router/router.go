package router

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gplume/todo-list/pkg/api"
	prome "github.com/gplume/todo-list/pkg/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Engine contains all gin-gonic settings and Routes
var Engine *gin.Engine

func init() {
	NewEngineAndRoutes()
}

// NewEngineAndRoutes returns gin.Engine with routes and options
func NewEngineAndRoutes() *gin.Engine {
	// Creates a router without any middleware by default
	r := gin.New()

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// group (will add /todo to all routes endpoints below:
	apiRoutes := r.Group("/todo")

	// statsMiddleWare() used by prometheus to stats requests latencies
	apiRoutes.Use(statsMiddleWare())

	// API endpoints (group suffix is added automatically)
	apiRoutes.GET("", api.ListTodos)
	apiRoutes.GET("/:id", api.GetTodo)
	apiRoutes.POST("", api.AddTodo)
	apiRoutes.PUT("", api.UpdateTodo)
	apiRoutes.DELETE("/:id", api.DeleteTodo)

	// Prometheus metrics by convention:
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	Engine = r
	return r
}

// StatsMiddleWare observe requests responses latencies on router Group (/todo) only
func statsMiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		code := strconv.Itoa(c.Writer.Status())
		elapsed := time.Since(start)
		msElapsed := elapsed / time.Millisecond
		prome.Vars.HTTPResponseLatencies.WithLabelValues(code, c.Request.Method).Observe(float64(msElapsed))
	}
}
