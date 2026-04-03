package admin

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersInfoUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	pagination := common.Pagination{Page: 1, PageSize: 10}

	tests := []struct {
		name    string
		roles   []string
		setup   func(m *MockUserRepository)
		wantErr bool
		err     error
	}{
		{
			name:  "Success With Roles",
			roles: []string{"ADMIN"},
			setup: func(m *MockUserRepository) {
				m.On("GetUsersInfo", ctx, []vo.Role{vo.AdminRole}, pagination).Return(&common.PaginatedResult[*entity.User]{}, nil)
			},
			wantErr: false,
		},
		{
			name:  "Success Without Roles (All Roles)",
			roles: []string{},
			setup: func(m *MockUserRepository) {
				m.On("GetUsersInfo", ctx, vo.AllRoles(), pagination).Return(&common.PaginatedResult[*entity.User]{}, nil)
			},
			wantErr: false,
		},
		{
			name:  "Invalid Role",
			roles: []string{"INVALID"},
			setup: func(m *MockUserRepository) {
				// No call expected
			},
			wantErr: true,
			err:     vo.ErrInvalidRole,
		},
		{
			name:  "Repository Error",
			roles: []string{"ADMIN"},
			setup: func(m *MockUserRepository) {
				m.On("GetUsersInfo", ctx, []vo.Role{vo.AdminRole}, pagination).Return(nil, assert.AnError)
			},
			wantErr: true,
			err:     assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockUserRepository)
			l := new(MockLogger)
			tt.setup(m)
			uc := NewGetUsersInfoUseCase(m, l)

			result, err := uc.Execute(ctx, pagination, tt.roles)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.err != nil {
					assert.ErrorIs(t, err, tt.err)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
			m.AssertExpectations(t)
		})
	}
}
