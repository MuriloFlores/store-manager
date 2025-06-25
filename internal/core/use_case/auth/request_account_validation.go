package auth

import (
	"context"
	"fmt"
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
	limiter        ports.RateLimiter
}

func NewRequestAccountValidationUseCase(
	userRepo repositories.UserRepository,
	tokenRepo repositories.ActionTokenRepository,
	tokenGenerator ports.SecureTokenGenerator,
	taskEnqueuer ports.TaskEnqueuer,
	logger ports.Logger,
	limiter ports.RateLimiter,
) *RequestAccountValidationUseCase {
	return &RequestAccountValidationUseCase{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		tokenGenerator: tokenGenerator,
		taskEnqueuer:   taskEnqueuer,
		logger:         logger,
		limiter:        limiter,
	}
}

func (uc *RequestAccountValidationUseCase) Execute(ctx context.Context, userEmail string) error {
	uc.logger.InfoLevel("Request account validation use case started", map[string]interface{}{"user_email": userEmail})

	user, err := uc.userRepo.FindByEmail(ctx, userEmail)
	if err != nil {
		uc.logger.ErrorLevel("Error finding user by ID", err, map[string]interface{}{"userEmail": userEmail})
		return err
	}

	rateLimitKey := fmt.Sprintf("rate-limit:resend-verification:%s", user.Email())
	limit := 5 * time.Minute

	allowed, err := uc.limiter.Allow(ctx, rateLimitKey, limit)
	if err != nil {
		uc.logger.ErrorLevel("Failed to check rate limit", err, map[string]interface{}{"user_id": user.ID()})
		return err
	}

	if !allowed {
		uc.logger.InfoLevel("Rate limit exceeded for resend verification email", map[string]interface{}{"user_id": user.ID(), "email": user.Email()})
		return &domain.ErrRateLimitExceeded{}
	}

	verificationTokenString, err := uc.tokenGenerator.Generate()
	if err != nil {
		uc.logger.ErrorLevel("Error generating verification token", err, nil)
		return err
	}

	actionToken := &domain.ActionToken{
		Token:     verificationTokenString,
		UserID:    user.ID(),
		Type:      domain.AccountVerification,
		ExpiresAt: time.Now().Add(time.Minute * 30),
	}

	if err = uc.tokenRepo.Create(ctx, actionToken); err != nil {
		uc.logger.ErrorLevel("Error creating action token", err, map[string]interface{}{"user_id": user.ID()})
		return err
	}

	jobData := &jobs.AccountVerificationJobData{
		UserName:         user.Name(),
		ToEmail:          user.Email(),
		VerificationLink: "https://app.muriloflores.xyz/auth/verify-account?token=" + actionToken.Token,
	}

	if err = uc.taskEnqueuer.EnqueueAccountVerification(jobData); err != nil {
		uc.logger.ErrorLevel("Error enqueuing account verification job", err, map[string]interface{}{"user_id": user.ID()})
		return err
	}

	uc.logger.InfoLevel("Account validation request processed successfully", map[string]interface{}{"user_id": user.ID(), "user_email": user.Email()})
	return nil
}
