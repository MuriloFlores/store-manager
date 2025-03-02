package main

import (
	"fmt"
	"log"
	"net/http"
	"store-manager/internal/application/services"
	"store-manager/internal/application/usecase/product"
	"store-manager/internal/application/usecase/raw_material"
	delivery "store-manager/internal/delivery/http"
	"store-manager/internal/delivery/router"
	"store-manager/internal/infrastructure/config"
	"store-manager/internal/infrastructure/persistence/connection"
	"store-manager/internal/infrastructure/persistence/orm"
	"store-manager/internal/infrastructure/persistence/repositories_gorm"
)

func init() {
	config.InitEnvConfig()
}

func main() {
	postgresDB := connection.ConnectPostgresDB()
	gormDB := orm.NewGormDB(postgresDB)

	productRepo := repositories_gorm.NewProductRepositoryGorm(gormDB)
	rawMaterialRepo := repositories_gorm.NewRawMaterialsRepositoryGorm(gormDB)

	createProductUseCase := product.NewCreateProductUseCase(productRepo)
	findProductByIdUseCase := product.NewFindProductByIdUseCase(productRepo)
	getAllProductsUseCase := product.NewGetAllProducts(productRepo)
	deleteProductsByIdUseCase := product.NewDeleteProductByIdUseCase(productRepo)
	updateProductUseCase := product.NewUpdateProductUseCase(productRepo)

	createRawMaterialUseCase := raw_material.NewCreateRawMaterialUseCase(rawMaterialRepo)
	findRawMaterialByIdUseCase := raw_material.NewFindRawMaterialByIdUseCase(rawMaterialRepo)
	getAllRawMaterialsUseCase := raw_material.NewGetAllRawMaterials(rawMaterialRepo)
	deleteRawMaterialUseCase := raw_material.NewDeleteRawMaterialUseCase(rawMaterialRepo)
	updateRawMaterialUseCase := raw_material.NewUpdateRawMaterialUseCase(rawMaterialRepo)

	productService := services.NewProductService(
		createProductUseCase,
		findProductByIdUseCase,
		getAllProductsUseCase,
		deleteProductsByIdUseCase,
		updateProductUseCase,
	)

	rawMaterialService := services.NewRawMaterialService(
		createRawMaterialUseCase,
		findRawMaterialByIdUseCase,
		getAllRawMaterialsUseCase,
		deleteRawMaterialUseCase,
		updateRawMaterialUseCase,
	)

	productHandler := delivery.NewProductHandler(productService)
	rawMaterialHandler := delivery.NewRawMaterialHandler(rawMaterialService)

	routing := router.ConfigureRoutes(productHandler, rawMaterialHandler)

	fmt.Println("Iniciando servidor....")
	log.Fatal(http.ListenAndServe(":8080", routing))
}
