package store

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/organization/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockStoreRepository struct {
	mock.Mock
}

func (m *MockStoreRepository) Save(ctx context.Context, store *entity.Store) error {
	args := m.Called(ctx, store)
	return args.Error(0)
}

type MockTenantProvisioner struct {
	mock.Mock
}

func (m *MockTenantProvisioner) CreateSchema(ctx context.Context, schemaName string) error {
	args := m.Called(ctx, schemaName)
	return args.Error(0)
}

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return fn(ctx)
}
