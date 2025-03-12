package services

import (
	"go.uber.org/zap"
	"store-manager/internal/application/DTOs/product_DTO"
	"store-manager/internal/application/usecase/product"
	"store-manager/internal/infrastructure/logging"
)

type productService struct {
	createProductUseCase     product.CreateProductUseCaseInterface
	findProductByIdUseCase   product.FindProductByIdUseCaseInterface
	getAllProductsUseCase    product.GetAllProductsUseCase
	deleteProductByIdUseCase product.DeleteProductByIdUseCaseInterface
	updateProductsUseCase    product.UpdateProductUseCaseInterface
}

type ProductServiceInterface interface {
	CreateProduct(input []product_DTO.CreateProductDTO) ([]product_DTO.ProductDTO, error)
	FindProductById(input []product_DTO.FindProductDTO) ([]product_DTO.ProductDTO, error)
	GetAllProducts() ([]product_DTO.ProductDTO, error)
	DeleteProductsByIds(input []product_DTO.FindProductDTO) error
	UpdateProducts(input []product_DTO.UpdateProductDTO) ([]product_DTO.ProductDTO, error)
}

func NewProductService(
	createProductUseCase product.CreateProductUseCaseInterface,
	findProductByIdUseCase product.FindProductByIdUseCaseInterface,
	getAllProductsUseCase product.GetAllProductsUseCase,
	deleteProductByIdUseCase product.DeleteProductByIdUseCaseInterface,
	updateProductsUseCase product.UpdateProductUseCaseInterface,
) ProductServiceInterface {
	return &productService{
		createProductUseCase:     createProductUseCase,
		findProductByIdUseCase:   findProductByIdUseCase,
		getAllProductsUseCase:    getAllProductsUseCase,
		deleteProductByIdUseCase: deleteProductByIdUseCase,
		updateProductsUseCase:    updateProductsUseCase,
	}
}

func (p *productService) CreateProduct(input []product_DTO.CreateProductDTO) ([]product_DTO.ProductDTO, error) {
	logging.Info("CreateProduct Journey", zap.String("Init", "CreateProductService"))
	return p.createProductUseCase.CreateProduct(input)
}

func (p *productService) FindProductById(input []product_DTO.FindProductDTO) ([]product_DTO.ProductDTO, error) {
	logging.Info("FindProduct Journey", zap.String("Init", "FindProductByIdService"))
	return p.findProductByIdUseCase.FindProductById(input)
}

func (p *productService) GetAllProducts() ([]product_DTO.ProductDTO, error) {
	logging.Info("GetAllProducts Journey", zap.String("Init", "GetAllProductsService"))
	return p.getAllProductsUseCase.GetAllProducts()
}

func (p *productService) DeleteProductsByIds(input []product_DTO.FindProductDTO) error {
	logging.Info("DeleteProductsByIds Journey", zap.String("Init", "DeleteProductsByIdsService"))
	return p.deleteProductByIdUseCase.DeleteProductById(input)
}

func (p *productService) UpdateProducts(input []product_DTO.UpdateProductDTO) ([]product_DTO.ProductDTO, error) {
	logging.Info("UpdateProducts Journey", zap.String("Init", "UpdateProductsService"))
	return p.updateProductsUseCase.UpdateProduct(input)
}
