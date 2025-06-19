package router

import (
	"github.com/gorilla/handlers" // Para o CORS
	"github.com/gorilla/mux"
	httphandlers "github.com/muriloFlores/StoreManager/infrastructure/web/http"
	"github.com/muriloFlores/StoreManager/infrastructure/web/middleware"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"log"
	"net/http"
)

func NewRouter(
	userHandler *httphandlers.UserHandler,
	authHandler *httphandlers.AuthHandler,
	tokenManager ports.TokenManager,
) http.Handler {
	r := mux.NewRouter()

	// --- Rotas Públicas ---
	r.HandleFunc("/create-user", userHandler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	// --- Rotas Protegidas ---
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(tokenManager))

	api.HandleFunc("/auth/change-password", authHandler.ChangePassword).Methods(http.MethodPut)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	log.Println("Roteador configurado com sucesso.")
	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)
}
