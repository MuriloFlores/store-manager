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
	itemHandler *httphandlers.ItemHandler,
	tokenManager ports.TokenManager,
) http.Handler {
	r := mux.NewRouter()

	// --- Rotas Públicas ---
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Connected to Order Manager API!"))
	})

	// --- Rota De Itens para todos os Usuários ---
	r.HandleFunc("/items", itemHandler.ListPublicItems).Methods(http.MethodGet)

	// --- Rota de Busca de Itens ---
	searchItemsHandler := http.HandlerFunc(itemHandler.SearchItem)
	r.Handle("/items/search/{param}", middleware.TryAuthMiddleware(tokenManager)(searchItemsHandler)).Methods(http.MethodGet)

	// --- Rotas de Autenticação e Criação de Usuários ---
	r.HandleFunc("/create-user", userHandler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)

	// --- Rotas para Validação de Conta ---
	r.HandleFunc("/verify-account", authHandler.ConfirmAccount).Methods(http.MethodGet)
	r.HandleFunc("/auth/resend-verification", authHandler.ResendVerificationEmail).Methods(http.MethodPost)

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
	api.HandleFunc("/user/{id}/role", userHandler.PromoteUser).Methods(http.MethodPatch)
	api.HandleFunc("/users", userHandler.ListUsers).Methods(http.MethodGet)
	api.HandleFunc("/users/search/{param}", userHandler.SearchUser).Methods(http.MethodGet)

	// --- Rotas de Itens ---
	api.HandleFunc("/items", itemHandler.ListInternalItems).Methods(http.MethodGet)
	api.HandleFunc("/item", itemHandler.CreateItem).Methods(http.MethodPost)
	api.HandleFunc("/item/{id}", itemHandler.DeleteItem).Methods(http.MethodDelete)
	api.HandleFunc("/item/{id}", itemHandler.FindItemByID).Methods(http.MethodGet)
	api.HandleFunc("/item/{sku}", itemHandler.FindItemBySKU).Methods(http.MethodGet)
	api.HandleFunc("/item/{id}", itemHandler.UpdateItem).Methods(http.MethodPatch)
	api.HandleFunc("/items/{id}/reactivate", itemHandler.ReactiveItem).Methods(http.MethodPost)

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	log.Println("Roteador configurado com sucesso.")
	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)
}
