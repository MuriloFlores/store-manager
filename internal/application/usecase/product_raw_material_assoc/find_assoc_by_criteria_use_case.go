package product_raw_material_assoc

import (
	"fmt"
	"go.uber.org/zap"
	"store-manager/internal/domain/entity"
	"store-manager/internal/domain/repositories"
	"store-manager/internal/infrastructure/logging"
)

type findAssocByCriteriaUseCase struct {
	assocRepo repositories.ProductAssocRawMaterialRepositoryInterface
}

type FindAssocByCriteriaUseCaseInterface interface {
	FindByCriteria(criteria string, info interface{}) ([]entity.ProductInterface, error)
}

func NewFindAssocByCriteriaUseCase(assocRepo repositories.ProductAssocRawMaterialRepositoryInterface) FindAssocByCriteriaUseCaseInterface {
	return &findAssocByCriteriaUseCase{
		assocRepo: assocRepo,
	}
}

func (uc *findAssocByCriteriaUseCase) FindByCriteria(criteria string, info interface{}) ([]entity.ProductInterface, error) {
	logging.Info("FindAssocByCriteria", zap.String("Init", "FindByCriteriaUseCase"))
	if criteria == "" {
		logging.Error("FindAssocByCriteria", zap.String("Error", "Criteria is empty"))
		return []entity.ProductInterface{}, nil
	}

	switch criteria {
	case "byProduct":
		slice := uc.convertInterfaceInStringSlice(info)
		return uc.findByProduct(slice)
	case "byRawMaterial":
		slice := uc.convertInterfaceInStringSlice(info)
		return uc.findByRawMaterial(slice)
	case "byActivated":
		status := info.(bool)
		return uc.findByActivated(status)
	default:

		logging.Info("FindAssocByCriteria Journey", zap.String("Finish", "FindAssocByCriteriaUseCase"))
		logging.Info("FindAssocByCriteria Journey", zap.String("Finish", "FindAssocByCriteriaService"))
		return []entity.ProductInterface{}, fmt.Errorf("invalid criteria")
	}
}

func (uc *findAssocByCriteriaUseCase) findByProduct(productId []string) ([]entity.ProductInterface, error) {
	logging.Info("findByProduct", zap.String("Init", "findByProduct"))
	products, err := uc.assocRepo.ListAssociationByProduct(productId)

	if err != nil {
		logging.Error("findByProduct", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, err
	}

	logging.Info("findByProduct", zap.String("Finish", "findByProduct"))
	logging.Info("FindAssocByCriteria Journey", zap.String("Finish", "FindAssocByCriteriaUseCase"))
	logging.Info("FindAssocByCriteria Journey", zap.String("Finish", "FindAssocByCriteriaService"))

	return products, nil
}

func (uc *findAssocByCriteriaUseCase) findByRawMaterial(rawMaterialId []string) ([]entity.ProductInterface, error) {
	logging.Info("findByRawMaterial", zap.String("Init", "findByRawMaterial"))
	products, err := uc.assocRepo.ListAssociationByRawMaterial(rawMaterialId)

	if err != nil {
		logging.Error("findByRawMaterial", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, err
	}

	logging.Info("findByRawMaterial", zap.String("Finish", "findByRawMaterial"))
	logging.Info("FindAssocByCriteria Journey", zap.String("Finish", "FindAssocByCriteriaUseCase"))
	logging.Info("FindAssocByCriteria Journey", zap.String("Finish", "FindAssocByCriteriaService"))

	return products, nil
}

func (uc *findAssocByCriteriaUseCase) findByActivated(activated bool) ([]entity.ProductInterface, error) {
	logging.Info("findByActivated", zap.String("Init", "findByActivated"))
	products, err := uc.assocRepo.ListAssociationByActivated(activated)

	if err != nil {
		logging.Error("findByActivated", zap.String("Error", err.Error()))
		return []entity.ProductInterface{}, err
	}

	logging.Info("findByActivated", zap.String("Finish", "findByActivated"))
	logging.Info("FindAssocByCriteria Journey", zap.String("Finish", "FindAssocByCriteriaUseCase"))
	logging.Info("FindAssocByCriteria Journey", zap.String("Finish", "FindAssocByCriteriaService"))

	return products, nil
}

func (uc *findAssocByCriteriaUseCase) convertInterfaceInStringSlice(info interface{}) []string {
	infoSlice, ok := info.([]interface{})
	if !ok {
		return []string{""}
	}

	var ids []string

	for _, value := range infoSlice {
		str, ok := value.(string)

		if !ok {
			return []string{""}
		}

		ids = append(ids, str)
	}

	return ids
}
