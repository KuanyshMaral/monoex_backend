package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"monoex_backend/database"
	"monoex_backend/internal/config"
	"monoex_backend/internal/handlers"
	"monoex_backend/internal/middleware"
	"monoex_backend/internal/repositories"
	"monoex_backend/internal/routes"
	"monoex_backend/internal/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	config *config.Config
	db     *sql.DB
	router *mux.Router
	server *http.Server

	// Repositories
	legislationRepo *repositories.LegislationRepository
	newsRepo        *repositories.NewsRepository
	reviewRepo      *repositories.ReviewRepository

	// Services
	legislationService *services.LegislationService
	newsService        *services.NewsService
	reviewService      *services.ReviewService
	adminService       *services.AdminService

	// Handlers
	legislationHandler *handlers.LegislationHandler
	newsHandler        *handlers.NewsHandler
	reviewHandler      *handlers.ReviewHandler
	adminHandler       *handlers.AdminHandler
}

func New() *App {
	return &App{}
}

func (a *App) Initialize() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.config = cfg

	// Initialize database (sql.DB) using DSN from config
	if err := a.initDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize migrations
	if err := database.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("✅ Migrations выполнены успешно")

	// Initialize repositories
	a.initRepositories()

	// Initialize services
	a.initServices()

	// Initialize handlers
	a.initHandlers()

	// Setup routes
	a.setupRoutes()

	// Setup server
	a.setupServer()

	return nil
}

func (a *App) initDB() error {
	db, err := sql.Open(a.config.DB.Driver, a.config.DB.DSN)
	if err != nil {
		return err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return err
	}

	// Configure connection pool
	db.SetMaxOpenConns(a.config.DB.MaxOpenConns)
	db.SetMaxIdleConns(a.config.DB.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(a.config.DB.ConnMaxLifetime) * time.Minute)

	a.db = db
	return nil
}

func (a *App) initRepositories() {
	a.legislationRepo = repositories.NewLegislationRepository(a.db)
	a.newsRepo = repositories.NewNewsRepository(a.db)
	a.reviewRepo = repositories.NewReviewRepository(a.db)
}

func (a *App) initServices() {
	a.legislationService = services.NewLegislationService(a.legislationRepo)
	a.newsService = services.NewNewsService(a.newsRepo)
	a.reviewService = services.NewReviewService(a.reviewRepo)
	a.adminService = services.NewAdminService(a.db)
}

func (a *App) initHandlers() {
	a.legislationHandler = handlers.NewLegislationHandler(a.legislationService)
	a.newsHandler = handlers.NewNewsHandler(a.newsService)
	a.reviewHandler = handlers.NewReviewHandler(a.reviewService)
	a.adminHandler = handlers.NewAdminHandler(a.adminService)
}

func (a *App) setupRoutes() {
	a.router = mux.NewRouter()

	// Apply global middleware
	a.router.Use(middleware.CORS())
	a.router.Use(middleware.Logging())
	a.router.Use(middleware.Recovery())

	// ✅ Health-check (можно открыть в браузере /health)
	a.router.HandleFunc("/health", a.healthCheck).Methods("GET")

	// ✅ Регистрация единственного админа (без middleware)
	a.router.HandleFunc("/register-admin", a.adminHandler.Register).Methods("POST")

	// ✅ Подключаем API роуты с админским middleware
	routes.RegisterAllRoutes(a.router, a.db, a.adminService)

	// ✅ Раздача загруженных файлов (статические файлы)
	uploadDir := "./uploads"
	a.router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))
}

func (a *App) setupServer() {
	a.server = &http.Server{
		Addr:         ":" + a.config.Server.Port,
		Handler:      a.router,
		ReadTimeout:  time.Duration(a.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(a.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(a.config.Server.IdleTimeout) * time.Second,
	}
}

func (a *App) Run() error {
	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", a.config.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	// Close database connection
	if err := a.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server exited")
	return nil
}

func (a *App) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}
