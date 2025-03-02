package repositories

import (
	"store-manager/internal/domain/entity"
)

type RawMaterialsRepositoryInterface interface {
	Save(rawMaterials []entity.RawMaterialInterface) error
	FindByIds(ids []string) ([]entity.RawMaterialInterface, error)
	Update(rawMaterials []entity.RawMaterialInterface) error
	DeleteByIds(ids []string) error
	GetAllRawMaterials() ([]entity.RawMaterialInterface, error)
	GetAllRawMaterialsByLimitRisk() ([]entity.RawMaterialInterface, error)
}
