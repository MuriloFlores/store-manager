package auth

import (
	"context"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/jobs"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"github.com/muriloFlores/StoreManager/internal/core/ports/repositories"
	"time"
)

type RequestAccountValidationUseCase struct {
	userRepo       repositories.UserRepository
	tokenRepo      repositories.ActionTokenRepository
	tokenGenerator ports.SecureTokenGenerator
	taskEnqueuer   ports.TaskEnqueuer
	logger         ports.Logger
}

func NewRequestAccountValidationUseCase(
	userRepo repositories.UserRepository,
	tokenRepo repositories.ActionTokenRepository,
	tokenGenerator ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	logger ports.Logger,
) *RequestAccountValidationUseCase {
	return &RequestAccountValidationUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		tokenGenerator: tokenGenerator,
		taskEnqueuer:   taskEnqueuer,
		logger:         logger,
	}
}

func (uc *RequestAccountValidationUseCase) Execute(ctx context.Context, userID string) error {
	uc.logger.InfoLevel("Request account validation use case started", map[string]interface{}{"user_id": userID})

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		uc.logger.ErrorLevel("Error finding user by ID", err, map[string]interface{}{"user_id": userID})
		return err
	}

	verificationTokenString, err := uc.tokenGenerator.Generate()
	if err != nil {
		uc.logger.ErrorLevel("Error generating verification token", err, nil)
		return err
	}

	actionToken := &domain.ActionToken{
		Token:     verificationTokenString,
		UserID:    userID,
		Type:      domain.AccountVerification,
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}

	if err = uc.tokenRepo.Create(ctx, actionToken); err != nil {
		uc.logger.ErrorLevel("Error creating action token", err, map[string]interface{}{"user_id": userID})
		return err
	}

	jobData := &jobs.AccountVerificationJobData{
		UserName:         user.Name(),
		ToEmail:          user.Email(),
		VerificationLink: "https://muriloflores.xyz/verify-account?token=" + actionToken.Token,
	}

	if err = uc.taskEnqueuer.EnqueueAccountVerification(jobData); err != nil {
		uc.logger.ErrorLevel("Error enqueuing account verification job", err, map[string]interface{}{"user_id": userID})
		return err
	}

	uc.logger.InfoLevel("Account validation request processed successfully", map[string]interface{}{"user_id": userID, "user_email": user.Email()})
	return nil
}
