package repository

import (
	"context"

	"github.com/MuriloFlores/order-manager/internal/_common"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/entity"
	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/infrastructure/database/model"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type userRepositoryImpl struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) ports.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Save(ctx context.Context, user *entity.User) error {
	userModel := model.ToModel(user)

	_, err := r.db.NewInsert().Model(userModel).Exec(ctx)
	return err
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	userModel := new(model.UserModel)

	err := r.db.NewSelect().
		Model(userModel).
		Where("email = ?", email).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return model.ToEntity(userModel)

}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	userModel := new(model.UserModel)

	err := r.db.NewSelect().
		Model(userModel).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return model.ToEntity(userModel)
}

func (r *userRepositoryImpl) GetUsersInfo(ctx context.Context, roles []vo.Role, pagination _common.Pagination) (*_common.PaginatedResult[*entity.User], error) {
	var userModels []model.UserModel

	roleStr := make([]string, 0, len(roles))
	for _, role := range roles {
		roleStr = append(roleStr, role.String())
	}

	allowedSorts := map[string]string{
		"created_at": "created_at",
		"username":   "username",
		"email":      "email",
		"roles":      "roles",
	}

	sortCol := pagination.Sort
	if val, ok := allowedSorts[sortCol]; ok {
		sortCol = val
	} else {
		sortCol = "name"
	}

	query := r.db.NewSelect().Model(&userModels)

	if pagination.Search != "" {
		searchTerm := "%" + pagination.Search + "%"
		query.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("username ILIKE ?", searchTerm).
				WhereOr("email ILIKE ?", searchTerm)
		})
	}

	if len(roleStr) > 0 {
		query.Where("roles ?| array[?]", bun.In(roleStr))
	}

	total, err := query.
		OrderExpr("? ?", bun.Ident(sortCol), bun.Safe(pagination.Direction)).
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset()).
		ScanAndCount(ctx)

	if err != nil {
		return nil, err
	}

	entities := make([]*entity.User, 0, len(userModels))
	for _, userModel := range userModels {
		enty, err := model.ToEntity(&userModel)
		if err != nil {
			return nil, err
		}

		entities = append(entities, enty)
	}

	result := _common.NewPaginatedResult(entities, int64(total), pagination)
	return result, nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	userModel := model.ToModel(user)

	_, err := r.db.NewUpdate().Model(userModel).Exec(ctx)
	return err
}
