package auth

import (
	"context"
	"time"

	"github.com/MuriloFlores/order-manager/internal/identity/domain/vo"
	"github.com/MuriloFlores/order-manager/internal/identity/ports"
	"github.com/MuriloFlores/order-manager/internal/identity/ports/auth"
)

type forgotPasswordUseCase struct {
	otpRepo             ports.OTPRepository
	userRepo            ports.UserRepository
	notificationService ports.NotificationService
	expiresIn           time.Duration
}

func NewForgotPassword(
	otpRepo ports.OTPRepository,
	userRepo ports.UserRepository,
	notificationService ports.NotificationService,
	expiresIn time.Duration,
) auth.ForgotPasswordUseCase {
	return &forgotPasswordUseCase{
		otpRepo:             otpRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
		expiresIn:           expiresIn,
	}
}

func (uc *forgotPasswordUseCase) Execute(ctx context.Context, email string) error {
	emailVO, err := vo.NewEmail(email)
	if err != nil {
		return err
	}

	user, err := uc.userRepo.FindByEmail(ctx, emailVO)
	if err != nil {
		//futuramente loggar o erro, nao retornar erro para o  http para evitar enumerations
		return nil
	}

	if user == nil {
		//futuramente logar o erro, nao retornar erro para o  http para evitar enumerations
		return nil
	}

	otpVO, err := vo.GenerateOTP()
	if err != nil {
		return err
	}

	err = uc.otpRepo.SaveOTP(ctx, user.Email(), otpVO, uc.expiresIn)
	if err != nil {
		return err
	}

	err = uc.notificationService.SendForgotPasswordEmail(ctx, user.Email(), otpVO)
	if err != nil {
		return err
	}

	return nil
}
