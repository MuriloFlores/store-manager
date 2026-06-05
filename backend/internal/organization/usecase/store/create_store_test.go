package store

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateStoreUseCase_Execute(t *testing.T) {
	ownerID := uuid.New()
	storeName := "Minha Loja"

	t.Run("should provision schema and save store in a transaction successfully", func(t *testing.T) {
		mockRepo := new(MockStoreRepository)
		mockProvisioner := new(MockTenantProvisioner)
		mockTx := new(MockTransactionManager)

		// Arrange: Sucesso total
		mockTx.On("Execute", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockProvisioner.On("CreateSchema", mock.Anything, mock.AnythingOfType("string")).Return(nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Store")).Return(nil)

		useCase := NewCreateStoreUseCase(mockRepo, mockProvisioner, mockTx)

		// Act
		err := useCase.Execute(context.Background(), storeName, ownerID)

		// Assert
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockProvisioner.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("should return error and abort transaction if provisioner fails", func(t *testing.T) {
		mockRepo := new(MockStoreRepository)
		mockProvisioner := new(MockTenantProvisioner)
		mockTx := new(MockTransactionManager)

		expectedErr := errors.New("failed to provision schema")

		// Arrange: Provisioner vai falhar
		// IMPORTANTE: O mockTx deve retornar nil para que a função anônima de transação seja invocada pelo mock.
		mockTx.On("Execute", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockProvisioner.On("CreateSchema", mock.Anything, mock.AnythingOfType("string")).Return(expectedErr)

		useCase := NewCreateStoreUseCase(mockRepo, mockProvisioner, mockTx)

		// Act
		err := useCase.Execute(context.Background(), storeName, ownerID)

		// Assert
		assert.ErrorIs(t, err, expectedErr)

		// Verifica se o provisioner foi chamado
		mockProvisioner.AssertExpectations(t)
		mockTx.AssertExpectations(t)

		// CRÍTICO: Se o provisioner falha, o repositório NUNCA deve ser chamado!
		mockRepo.AssertNotCalled(t, "Save")
	})

	t.Run("should return error if repository fails to save", func(t *testing.T) {
		mockRepo := new(MockStoreRepository)
		mockProvisioner := new(MockTenantProvisioner)
		mockTx := new(MockTransactionManager)

		expectedErr := errors.New("database connection lost")

		// Arrange: Provisioner funciona, mas Repository falha
		// IMPORTANTE: O mockTx deve retornar nil para que a função anônima de transação seja invocada pelo mock.
		mockTx.On("Execute", mock.Anything, mock.AnythingOfType("func(context.Context) error")).Return(nil)
		mockProvisioner.On("CreateSchema", mock.Anything, mock.AnythingOfType("string")).Return(nil)
		mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*entity.Store")).Return(expectedErr)
		useCase := NewCreateStoreUseCase(mockRepo, mockProvisioner, mockTx)

		// Act
		err := useCase.Execute(context.Background(), storeName, ownerID)

		// Assert
		assert.ErrorIs(t, err, expectedErr)

		// Ambos devem ter sido chamados antes de a transação falhar
		mockProvisioner.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})
}
