package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/qulaz/gas-price-test/internal/config"
	"github.com/qulaz/gas-price-test/internal/entity"
	"github.com/qulaz/gas-price-test/internal/usecase"
	v1 "github.com/qulaz/gas-price-test/internal/usecase/controller/http/v1"
	"github.com/qulaz/gas-price-test/internal/usecase/repository"
	"github.com/qulaz/gas-price-test/pkg/cache/memory"
	"github.com/qulaz/gas-price-test/pkg/logging"
	"github.com/qulaz/gas-price-test/pkg/shutdown"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	loggerConfig := logging.LoggerConfig{Mode: logging.ProductionMode, Level: logging.InfoLevel}
	if cfg.API.Debug {
		loggerConfig = logging.LoggerConfig{
			Mode: logging.DevelopmentMode, Level: logging.DebugLevel,
		}
	}

	logger, err := logging.NewZapLogger(loggerConfig)
	if err != nil {
		panic(err)
	}

	repo := repository.NewHttpGasTransactionRepo(cfg.API.TransactionsUrl)

	gasGraphCache := memory.NewExpiredCache[string, entity.GasGraphResult]()

	gasUseCase := usecase.NewGasGraphUseCase(
		repo,
		gasGraphCache,
		cfg.API.GasGraphTtl,
		logger,
	)

	router := gin.New()
	v1.NewRouter(router, cfg, logger, gasUseCase)

	server := &http.Server{ //nolint: exhaustruct
		Addr:         fmt.Sprintf("%s:%s", cfg.API.Host, cfg.API.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go shutdown.Graceful(
		[]os.Signal{
			syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM,
		},
		server,
	)

	logger.Infow(fmt.Sprintf("🚀 Starting server at http://%s:%s", cfg.API.Host, cfg.API.Port))

	if err = server.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Infow("✅ Server shutdown successfully")
		default:
			logger.Errorw("🔥 Server stopped due error", "error", err.Error())
		}
	}
}
