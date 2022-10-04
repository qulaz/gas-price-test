package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/qulaz/gas-price-test/docs"
	"github.com/qulaz/gas-price-test/internal/config"
	"github.com/qulaz/gas-price-test/internal/usecase"
	"github.com/qulaz/gas-price-test/pkg/logging"
)

// NewRouter -.
// Swagger spec:
// @title       Gas Price Test Task
// @version     1.0
// @BasePath    /v1
func NewRouter(handler *gin.Engine, cfg *config.Config, l logging.ContextLogger, t usecase.GasGraph) {
	handler.Use(cors.Default())
	handler.Use(logging.GinLogging(l))
	handler.Use(gin.Recovery())

	if cfg.API != nil {
		fmt.Println(cfg.API.Domain)
		if cfg.API.Domain != "" {
			docs.SwaggerInfo.Host = cfg.API.Domain
			docs.SwaggerInfo.Schemes = []string{"https"}
		} else {
			docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", cfg.API.Port)
		}
	}

	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	h := handler.Group("/v1")
	{
		newGasGraphRoutes(h, t, l)
	}
}
