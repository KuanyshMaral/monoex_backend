package routes

import (
	"database/sql"
	"github.com/gorilla/mux"
	"monoex_backend/internal/handlers"
	"monoex_backend/internal/middleware"
	"monoex_backend/internal/repositories"
	"monoex_backend/internal/services"
)

// adminService нужен для middleware, чтобы проверять админа
func RegisterAllRoutes(r *mux.Router, db *sql.DB, adminService *services.AdminService) {

	legRepo := repositories.NewLegislationRepository(db)
	legService := services.NewLegislationService(legRepo)
	legHandler := handlers.NewLegislationHandler(legService)

	newsRepo := repositories.NewNewsRepository(db)
	newsService := services.NewNewsService(newsRepo)
	newsHandler := handlers.NewNewsHandler(newsService)

	reviewRepo := repositories.NewReviewRepository(db)
	reviewService := services.NewReviewService(reviewRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	// --- Middleware для админа ---
	adminMiddleware := middleware.AdminMiddleware(adminService)

	adminHandler := handlers.NewAdminHandler(adminService)
	// Регистрация единственного админа
	r.HandleFunc("/admin/register", adminHandler.Register).Methods("POST")

	// --- CRUD законы (только для админа) ---
	r.Handle("/legislations", adminMiddleware(legHandler.Create)).Methods("POST")
	r.Handle("/legislations", adminMiddleware(legHandler.GetAll)).Methods("GET")
	r.Handle("/legislations/{id}", adminMiddleware(legHandler.GetByID)).Methods("GET")
	r.Handle("/legislations/{id}", adminMiddleware(legHandler.Update)).Methods("PUT", "PATCH")
	r.Handle("/legislations/{id}", adminMiddleware(legHandler.Delete)).Methods("DELETE")
	r.Handle("/files/legislations", adminMiddleware(legHandler.UploadFile)).Methods("POST")

	// --- CRUD новости (только для админа, кроме публичных) ---
	r.Handle("/news", adminMiddleware(newsHandler.Create)).Methods("POST")
	r.Handle("/news", adminMiddleware(newsHandler.GetAll)).Methods("GET")
	r.Handle("/news/{id:[0-9]+}", adminMiddleware(newsHandler.GetByID)).Methods("GET")
	r.Handle("/news/{id:[0-9]+}", adminMiddleware(newsHandler.Update)).Methods("PUT", "PATCH")
	r.Handle("/news/{id:[0-9]+}", adminMiddleware(newsHandler.Delete)).Methods("DELETE")

	// --- Получить по ссылке (public) ---
	r.HandleFunc("/news/by-link/{link}", newsHandler.GetByLink).Methods("GET")

	// --- Получить только опубликованные новости (public) ---
	r.HandleFunc("/news/published", newsHandler.GetPublished).Methods("GET")

	// --- Publish / Unpublish (только админ) ---
	r.Handle("/news/{id:[0-9]+}/publish", adminMiddleware(newsHandler.Publish)).Methods("POST")
	r.Handle("/news/{id:[0-9]+}/unpublish", adminMiddleware(newsHandler.Unpublish)).Methods("POST")

	// --- Загрузка изображения (только админ) ---
	r.Handle("/files/news-image", adminMiddleware(newsHandler.UploadImage)).Methods("POST")

	// --- CRUD отзывы (только админ) ---
	r.Handle("/reviews", adminMiddleware(reviewHandler.Create)).Methods("POST")
	r.Handle("/reviews", adminMiddleware(reviewHandler.GetAll)).Methods("GET")
	r.Handle("/reviews/{id:[0-9]+}", adminMiddleware(reviewHandler.GetByID)).Methods("GET")
	r.Handle("/reviews/{id:[0-9]+}", adminMiddleware(reviewHandler.Update)).Methods("PUT", "PATCH")
	r.Handle("/reviews/{id:[0-9]+}", adminMiddleware(reviewHandler.Delete)).Methods("DELETE")
	r.Handle("/files/reviews", adminMiddleware(reviewHandler.UploadFile)).Methods("POST")
}
