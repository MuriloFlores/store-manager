package repositories

import "store-manager/internal/domain/entity"

type ProductRepositoryInterface interface {
	Save(product []entity.ProductInterface) error
	GetAllProducts() ([]entity.ProductInterface, error)
	FindByIds(ids []string) ([]entity.ProductInterface, error)
	Update(products []entity.ProductInterface) error
	DeleteByIds(ids []string) error
}
