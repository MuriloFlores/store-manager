package store

import (
	"context"
	"fmt"

	"github.com/MuriloFlores/order-manager/internal/organization/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/organization/ports"
	"github.com/google/uuid"
)

type CreateStoreUseCase struct {
	storeRepo          ports.StoreRepository
	tenantProvisioner  ports.TenantProvisioner
	transactionManager ports.TransactionManager
}

func NewCreateStoreUseCase(
	storeRepo ports.StoreRepository,
	tenantProvisioner ports.TenantProvisioner,
	transactionManager ports.TransactionManager,
) *CreateStoreUseCase {
	return &CreateStoreUseCase{
		storeRepo:          storeRepo,
		tenantProvisioner:  tenantProvisioner,
		transactionManager: transactionManager,
	}
}

func (uc *CreateStoreUseCase) Execute(ctx context.Context, storeName string, ownerID uuid.UUID) error {
	store, err := entity.NewStore(storeName, ownerID)
	if err != nil {
		return fmt.Errorf("failed to create store entity: %w", err)
	}

	err = uc.transactionManager.Execute(ctx, func(txCtx context.Context) error {
		if err := uc.tenantProvisioner.CreateSchema(txCtx, store.SchemaName.String()); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}

		if err := uc.storeRepo.Save(txCtx, store); err != nil {
			return fmt.Errorf("failed to save store: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	return nil
}
