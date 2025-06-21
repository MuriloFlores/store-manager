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
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("01101000 01101001 01101100 01101100 00100000 01100100 01100101 01100001 01101110 00100000 11101001 00100000 01110110 01101001 01100001 01100100 01101111 "))
	})

	// --- Rotas de Autenticação e Criação de Usuários ---
	r.HandleFunc("/create-user", userHandler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	// --- Rotas para Validação de Conta ---
	r.HandleFunc("/verify-account", authHandler.ConfirmAccount).Methods(http.MethodGet)

	// --- Rotas Protegidas ---
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(tokenManager))

	// --- Rotas para alteração de Senha e Email ---
	api.HandleFunc("/auth/change-password", authHandler.ChangePassword).Methods(http.MethodPut)
	api.HandleFunc("/auth/confirm-email", authHandler.ConfirmEmail).Methods(http.MethodPost)

	// --- Rotas de CRUD de Usuários ---
	api.HandleFunc("/user/{id}", userHandler.DeleteUser).Methods(http.MethodDelete)
	api.HandleFunc("/user/{id}", userHandler.FindUserByID).Methods(http.MethodGet)
	api.HandleFunc("/user", userHandler.FindUserByEmail).Methods(http.MethodGet).Queries("email", "{email}")
	api.HandleFunc("/user/{id}", userHandler.UpdateUser).Methods(http.MethodPut)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	log.Println("Roteador configurado com sucesso.")
	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)
}
