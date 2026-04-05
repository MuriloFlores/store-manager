package user

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/dto"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUserUseCase_Execute(t *testing.T) {
	pepper := "test-pepper"

	tests := []struct {
		name    string
		input   dto.CreateUserInput
		setup   func(mockRepo *MockUserRepository, mockTx *MockTransactionManager)
		wantErr bool
	}{
		{
			name: "Success",
			input: dto.CreateUserInput{
				Username: "murilo",
				Email:    "murilo@test.com",
				Password: "StrongPass123!",
				Roles:    []string{"ADMIN"},
			},
			setup: func(mockRepo *MockUserRepository, mockTx *MockTransactionManager) {
				mockTx.On("Execute", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
					return u.Username() == "murilo" && u.Email().String() == "murilo@test.com"
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid Email",
			input: dto.CreateUserInput{
				Username: "murilo",
				Email:    "invalid-email",
				Password: "StrongPass123!",
				Roles:    []string{"ADMIN"},
			},
			setup: func(mockRepo *MockUserRepository, mockTx *MockTransactionManager) {},
			wantErr: true,
		},
		{
			name: "Transaction Start Error",
			input: dto.CreateUserInput{
				Username: "murilo",
				Email:    "murilo@test.com",
				Password: "StrongPass123!",
				Roles:    []string{"ADMIN"},
			},
			setup: func(mockRepo *MockUserRepository, mockTx *MockTransactionManager) {
				mockTx.On("Execute", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "Repository Save Error (Rollback Simulation)",
			input: dto.CreateUserInput{
				Username: "murilo",
				Email:    "murilo@test.com",
				Password: "StrongPass123!",
				Roles:    []string{"ADMIN"},
			},
			setup: func(mockRepo *MockUserRepository, mockTx *MockTransactionManager) {
				mockTx.On("Execute", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("Save", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockLogger := new(MockLogger)
			mockTx := new(MockTransactionManager)
			
			tt.setup(mockRepo, mockTx)

			uc := NewCreateUserService(mockRepo, mockLogger, mockTx, pepper)
			err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				mockRepo.AssertExpectations(t)
				mockTx.AssertExpectations(t)
			}
		})
	}
}
