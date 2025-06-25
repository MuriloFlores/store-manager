package user

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/jobs"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"github.com/muriloFlores/StoreManager/internal/core/value_objects"
)

type PromoteUserUseCase struct {
	userRepo repositories.UserRepository
	logger   ports.Logger
	enqueuer ports.TaskEnqueuer
}

func NewPromoteUseCase(userRepo repositories.UserRepository, logger ports.Logger, enqueuer ports.TaskEnqueuer) *PromoteUserUseCase {
	return &PromoteUserUseCase{
		userRepo: userRepo,
		logger:   logger,
		enqueuer: enqueuer,
	}
}

func (uc *PromoteUserUseCase) Execute(ctx context.Context, actor *domain.Identity, targetId, newRole string) (*domain.User, error) {
	uc.logger.InfoLevel("Promoting user", map[string]interface{}{"actor": actor.UserID, "targetId": targetId, "newRole": newRole})

	targetUser, err := uc.userRepo.FindByID(ctx, targetId)
	if err != nil {
		uc.logger.ErrorLevel("Target user not found", err, map[string]interface{}{"targetId": targetId})
		return nil, err
	}

	if err = uc.canPromote(actor, targetUser, value_objects.Role(newRole)); err != nil {
		uc.logger.ErrorLevel("Permission denied for promotion", err, map[string]interface{}{"actor": actor.UserID, "targetId": targetId, "newRole": newRole})
		return nil, err
	}

	if err = targetUser.ChangeRole(value_objects.Role(newRole)); err != nil {
		uc.logger.ErrorLevel("Failed to change user role", err, map[string]interface{}{"targetId": targetId, "newRole": newRole})
		return nil, err
	}

	if err = uc.userRepo.Update(ctx, targetUser); err != nil {
		uc.logger.ErrorLevel("Failed to update user role in repository", err, map[string]interface{}{"targetId": targetId, "newRole": newRole})
		return nil, err
	}

	uc.logger.InfoLevel("User role promoted successfully", map[string]interface{}{"targetId": targetId, "newRole": newRole})

	jobData := &jobs.PromotionNotificationJobData{
		UserName: targetUser.Name(),
		ToEmail:  targetUser.Email(),
		NewRole:  newRole,
	}

	if err = uc.enqueuer.EnqueuePromotionNotification(jobData); err != nil {
		uc.logger.ErrorLevel("Failed to enqueue promotion notification", err, map[string]interface{}{"targetId": targetId, "newRole": newRole})
		return nil, err
	}

	uc.logger.InfoLevel("Promotion notification enqueued successfully", map[string]interface{}{"targetId": targetId, "newRole": newRole})

	return targetUser, nil
}

func (uc *PromoteUserUseCase) canPromote(actor *domain.Identity, targetUser *domain.User, newRole value_objects.Role) error {
	actorRole := actor.Role

	switch actorRole {
	case value_objects.Admin:
		return nil
	case value_objects.Manager:
		if targetUser.Role() == value_objects.Admin || targetUser.Role() == value_objects.Manager {
			return &domain.ErrForbidden{Action: "manager cannot change role of other managers or admins"}
		}

		if newRole == value_objects.Admin || newRole == value_objects.Manager {
			return &domain.ErrForbidden{Action: "manager cannot promote to admin or manager"}
		}

		return nil
	}

	return &domain.ErrForbidden{Action: "user does not have permission to promote"}
}
