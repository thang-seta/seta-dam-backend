package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/user/seta-dam-backend/services/core-go/internal/common"
	"github.com/user/seta-dam-backend/services/core-go/internal/config"
	"github.com/user/seta-dam-backend/services/core-go/internal/database"
	"github.com/user/seta-dam-backend/services/core-go/pkg/logger"

	folderRepo "github.com/user/seta-dam-backend/services/core-go/internal/folder/repository"
	folderTransport "github.com/user/seta-dam-backend/services/core-go/internal/folder/transport"
	folderUsecase "github.com/user/seta-dam-backend/services/core-go/internal/folder/usecase"

	metaRepo "github.com/user/seta-dam-backend/services/core-go/internal/metadata/repository"
	metaTransport "github.com/user/seta-dam-backend/services/core-go/internal/metadata/transport"
	metaUsecase "github.com/user/seta-dam-backend/services/core-go/internal/metadata/usecase"

	permRepo "github.com/user/seta-dam-backend/services/core-go/internal/permission/repository"
	permTransport "github.com/user/seta-dam-backend/services/core-go/internal/permission/transport"
	permUsecase "github.com/user/seta-dam-backend/services/core-go/internal/permission/usecase"
)

func main() {
	// 1. Initialize Logger
	logger.Init()
	logger.Log.Info("Starting core-go service...")

	// 2. Load Config
	cfg := config.Load()

	// 3. Connect Database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Log.Info("Database connection established")

	// 4. Instantiate Layers
	// Repositories
	permRepository := permRepo.NewPostgresRepository(db)
	folderRepository := folderRepo.NewPostgresRepository(db)
	metaRepository := metaRepo.NewPostgresRepository(db)

	// Usecases
	permEvaluator := permUsecase.NewPermissionEvaluator(permRepository)
	folderUC := folderUsecase.NewFolderUsecase(folderRepository, permEvaluator)
	metaUC := metaUsecase.NewMetadataUsecase(metaRepository, folderRepository, permEvaluator)

	// Transport Handlers
	folderHandler := folderTransport.NewFolderHandler(folderUC)
	metaHandler := metaTransport.NewMetadataHandler(metaUC)
	permHandler := permTransport.NewPermissionHandler(permEvaluator)

	// 5. Router Setup
	r := chi.NewRouter()

	// Standard chi middlewares
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	// Custom user context parsing middleware
	r.Use(common.ExtractUserContext)

	// Health Check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"OK"}`))
	})

	// API Routes (registered with base prefix `/api`)
	r.Route("/api", func(apiRouter chi.Router) {
		folderHandler.RegisterRoutes(apiRouter)
		metaHandler.RegisterRoutes(apiRouter)
		permHandler.RegisterRoutes(apiRouter)
	})

	// 6. Graceful Shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		logger.Log.Info(fmt.Sprintf("Server listening on port %s", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown", "error", err)
	}

	logger.Log.Info("Server stopped")
}
