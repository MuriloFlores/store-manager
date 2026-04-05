package auth

import (
	"context"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/security"
)

type forgotPasswordUseCase struct {
	otpRepo             ports.OTPRepository
	userRepo            ports.UserRepository
	notificationService ports.NotificationService
	logger              ports.Logger
	expiresIn           time.Duration
}

func NewForgotPassword(
	otpRepo ports.OTPRepository,
	userRepo ports.UserRepository,
	notificationService ports.NotificationService,
	logger ports.Logger,
	expiresIn time.Duration,
) security.ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		otpRepo:             otpRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
		logger:              logger,
		expiresIn:           expiresIn,
	}
}

func (uc *forgotPasswordUseCase) Execute(ctx context.Context, email string) error {
	uc.logger.Debug("starting forgot password flow", "email", email)

	emailVO, err := vo.NewEmail(email)
	if err != nil {
		uc.logger.Info("invalid email provided in forgot password", "email", email)
		return err
	}

	user, err := uc.userRepo.FindByEmail(ctx, emailVO)
	if err != nil {
		uc.logger.Error("error finding user by email", err, "email", email)
		//futuramente loggar o erro, nao retornar erro para o  http para evitar enumerations
		return nil
	}

	if user == nil {
		uc.logger.Info("user not found in forgot password", "email", email)
		//futuramente logar o erro, nao retornar erro para o  http para evitar enumerations
		return nil
	}

	otpVO, err := vo.GenerateOTP()
	if err != nil {
		uc.logger.Error("failed to generate OTP", err, "email", email)
		return err
	}

	err = uc.otpRepo.SaveOTP(ctx, user.Email(), otpVO, uc.expiresIn)
	if err != nil {
		uc.logger.Error("failed to save OTP", err, "email", email)
		return err
	}

	err = uc.notificationService.SendForgotPasswordEmail(ctx, user.Email(), otpVO)
	if err != nil {
		uc.logger.Error("failed to send forgot password email", err, "email", email)
		return err
	}

	uc.logger.Info("forgot password OTP sent successfully", "email", email)
	return nil
}
