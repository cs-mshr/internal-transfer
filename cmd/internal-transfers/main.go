package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/chandra-shekhar/internal-transfers/internal/config"
	"github.com/chandra-shekhar/internal-transfers/internal/database"
	"github.com/chandra-shekhar/internal-transfers/internal/handler"
	"github.com/chandra-shekhar/internal-transfers/internal/logger"
	"github.com/chandra-shekhar/internal-transfers/internal/repository"
	"github.com/chandra-shekhar/internal-transfers/internal/router"
	"github.com/chandra-shekhar/internal-transfers/internal/server"
	"github.com/chandra-shekhar/internal-transfers/internal/service"
	"github.com/joho/godotenv"
)

const DefaultContextTimeout = 30

func main() {
	// Load .env file explicitly
	err := godotenv.Load()
	if err != nil {
		// It's ok if .env doesn't exist
		println("No .env file found, using environment variables")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Primary.Env)

	// Run migrations in non-local environments
	if cfg.Primary.Env != "local" {
		if err := database.Migrate(context.Background(), &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	// Initialize server
	srv, err := server.New(cfg, &log)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	// Initialize repositories, services, and handlers
	repos := repository.NewRepositories(srv)
	services := service.NewServices(srv, repos)
	handlers := handler.NewHandlers(srv, services)

	// Initialize router
	r := router.NewRouter(srv, handlers, services)

	// Setup HTTP server
	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Start server
	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}
	stop()
	cancel()

	log.Info().Msg("server exited properly")
}
