package admin

import (
	"context"
	"testing"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/stretchr/testify/assert"
)

func TestGetUsersInfoUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	pagination := _common.Pagination{Page: 1, PageSize: 10}

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
				m.On("GetUsersInfo", ctx, []vo.Role{vo.AdminRole}, pagination).Return(&_common.PaginatedResult[*entity.User]{}, nil)
			},
			wantErr: false,
		},
		{
			name:  "Success Without Roles (All Roles)",
			roles: []string{},
			setup: func(m *MockUserRepository) {
				m.On("GetUsersInfo", ctx, vo.AllRoles(), pagination).Return(&_common.PaginatedResult[*entity.User]{}, nil)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := new(MockUserRepository)
			tt.setup(m)
			uc := NewGetUsersInfoUseCase(m)

			result, err := uc.Execute(ctx, pagination, tt.roles)

			if tt.wantErr {
				assert.Error(t)
				assert.Nil(t, result)
				if tt.err != nil {
					assert.Equal(t, tt.err, err)
				}
			} else {
				assert.NoError(t)
				assert.NotNil(t, result)
			}
			m.AssertExpectations(t)
		})
	}
}
