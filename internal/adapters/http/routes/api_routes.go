package routes

import (
	"go-chi-boilerplate/internal/adapters/database/postgresql"
	"go-chi-boilerplate/internal/adapters/http/handlers"
	"go-chi-boilerplate/internal/app/usecases/auth"
	"go-chi-boilerplate/internal/app/usecases/user"
	"go-chi-boilerplate/internal/core/services"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

// AddAPIRoutes registers all API v1 routes
func AddAPIRoutes(rg chi.Router, db *postgresql.PostgresDB, logger *slog.Logger, jwtSecret string) {
	// Initialize repos, services, usecases
	userRepo := postgresql.NewUserRepository(db.DB)
	userUsecase := user.NewUserUsecase(userRepo)

	jwtService := services.NewJWTService(jwtSecret)
	authUsecase := auth.NewAuthUsecase(userRepo, jwtService)

	api := chi.NewRouter()

	// Versioned API prefix
	v1 := chi.NewRouter()

	// Auth routes
	v1.Post("/auth/login", handlers.LoginHandler(authUsecase))

	// Protected routes
	v1.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtService.TokenAuth))
		r.Use(jwtauth.Authenticator(jwtService.TokenAuth))
		r.Get("/profile", handlers.GetProfile(userUsecase))
		// User management routes
		r.Post("/users", handlers.CreateUserHandler(userUsecase))
		r.Get("/users", handlers.ListUsersHandler(userUsecase))
		r.Get("/users/{id}", handlers.GetUserHandler(userUsecase))
		r.Put("/users/{id}", handlers.UpdateUserHandler(userUsecase))
		r.Patch("/users/{id}/status", handlers.ChangeUserStatusHandler(userUsecase))
		r.Delete("/users/{id}", handlers.DeleteUserHandler(userUsecase))
	})

	// Mount versioned API
	api.Mount("/v1", v1)

	// Attach API router to main router
	rg.Mount("/api", api)
}
