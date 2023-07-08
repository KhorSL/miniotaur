package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/khorsl/minio_tutorial/api/v1/router"
	lgg "github.com/khorsl/minio_tutorial/common/log/logger"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := lgg.NewLoggerWrapper(os.Getenv("DEFAULT_LOGGER_TYPE"), context.Background())

	r := router.Initialize()
	router.ListRoutes(r)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		logger.Info("Starting server on port 8080...", nil)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-signalChan
	logger.Info("Shutting down server...", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown", map[string]interface{}{
			"error": err,
		})
	}

	logger.Info("Server gracefully stopped.", nil)
}
