package product_raw_material_assoc

import (
	"go.uber.org/zap"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type findAllAssocUseCase struct {
	assocRepo repositories.ProductAssocRawMaterialRepositoryInterface
}

type FindAllAssociationsUseCaseInterface interface {
	GetAllAssociations() ([]entity.ProductInterface, error)
}

func NewGetAllAssocUseCase(assocRepo repositories.ProductAssocRawMaterialRepositoryInterface) FindAllAssociationsUseCaseInterface {
	return &findAllAssocUseCase{
		assocRepo: assocRepo,
	}
}

func (uc *findAllAssocUseCase) GetAllAssociations() ([]entity.ProductInterface, error) {
	logging.Info("GetAllAssoc Journey", zap.String("Init", "GetAllAssocUseCase"))

	products, err := uc.assocRepo.GetAllAssociations()
	if err != nil {
		logging.Error("GetAllAssoc Journey", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, err
	}

	logging.Info("GetAllAssoc Journey", zap.String("Finish", "GetAllAssocUseCase"))
	logging.Info("GetAllAssoc Journey", zap.String("Finish", "GetAllAssocService"))

	return products, nil
}
