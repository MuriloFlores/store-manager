package product_raw_material_assoc

import (
	"errors"
	"go.uber.org/zap"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type deleteAssocUseCase struct {
	assocRepo repositories.ProductAssocRawMaterialRepositoryInterface
}

type DeleteAssocUseCaseInterface interface {
	DeleteByIds(productIds, materialsIds []string) error
}

func NewDeleteAssocUseCase(assocRepo repositories.ProductAssocRawMaterialRepositoryInterface) DeleteAssocUseCaseInterface {
	return &deleteAssocUseCase{
		assocRepo: assocRepo,
	}
}

func (uc *deleteAssocUseCase) DeleteByIds(productIds, materialsIds []string) error {
	logging.Info("DeleteAssoc Journey", zap.String("Init", "DeleteByIdsUseCase"))
	if len(productIds) == 0 || len(materialsIds) == 0 {
		logging.Error("DeleteAssoc Journey", zap.String("Error", "Invalid ids"))
		return errors.New("invalid ids")
	}

	err := uc.assocRepo.DeleteAssociation(productIds, materialsIds)
	if err != nil {
		logging.Error("DeleteAssoc Journey", zap.String("Error", err.Error()))
		return err
	}

	logging.Info("DeleteAssoc Journey", zap.String("Finish", "DeleteByIdsUseCase"))
	logging.Info("DeleteAssoc Journey", zap.String("Finish", "DeleteAssocService"))
	return nil
}
