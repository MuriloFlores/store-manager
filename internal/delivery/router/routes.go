package router

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	delivery "store-manager/internal/delivery/http"

	_ "store-manager/docs"
)

// @title Store Manager API
// @version 1.0
// @description API para gerenciar lojas.
// @host localhost:8080
// @BasePath /
func ConfigureRoutes(productHandler delivery.ProductHandlerInterface, rawMaterialHandler delivery.RawMaterialHandlerInterface) http.Handler {
	router := mux.NewRouter()

	// Rotas de Produtos
	router.HandleFunc("/products/insert", productHandler.CreateProduct).Methods("POST")
	router.HandleFunc("/products/get-by-ids", productHandler.FindProductById).Methods("GET")
	router.HandleFunc("/products/get-all", productHandler.GetAllProducts).Methods("GET")
	router.HandleFunc("/products/delete-by-ids", productHandler.DeleteProductsByIds).Methods("DELETE")
	router.HandleFunc("/products/update", productHandler.UpdateProduct).Methods("PUT")

	// Rotas de Matéria-Prima
	router.HandleFunc("/raw-material/insert", rawMaterialHandler.CreateRawMaterial).Methods("POST")
	router.HandleFunc("/raw-material/get-by-ids", rawMaterialHandler.FindRawMaterial).Methods("GET")
	router.HandleFunc("/raw-material/get-all", rawMaterialHandler.GetAllRawMaterials).Methods("GET")
	router.HandleFunc("/raw-material/delete-by-ids", rawMaterialHandler.DeleteRawMaterial).Methods("DELETE")
	router.HandleFunc("/raw-material/update", rawMaterialHandler.UpdateRawMaterial).Methods("PUT")

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return applyCors(router)
}

func applyCors(h http.Handler) http.Handler {
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:5173"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Origin", "Content-Length", "Content-Type", "Authorization"})
	allowCredentials := handlers.AllowCredentials()
	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders, allowCredentials)(h)
}
