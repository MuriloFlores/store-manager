package email

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
	"github.com/muriloFlores/StoreManager/internal/core/domain/jobs"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
)

type EmailProcessor struct {
	sender          ports.NotificationSender
	templateManager ports.TemplateManager
}

func NewEmailProcessor(sender ports.NotificationSender, templateManager ports.TemplateManager) ports.Worker {
	return &EmailProcessor{
		sender:          sender,
		templateManager: templateManager,
	}
}

func (p *EmailProcessor) HandlePasswordReset(ctx context.Context, taskPayload []byte) error {
	var payload jobs.PasswordResetJobData
	if err := json.Unmarshal(taskPayload, &payload); err != nil {
		return fmt.Errorf("unmarshal email password_reset payload: %w", err)
	}

	templateData := map[string]interface{}{
		"UserName":  payload.UserName,
		"ResetLink": payload.ResetLink,
	}

	htmlBody, err := p.templateManager.Render("password_reset.html", templateData)
	if err != nil {
		return fmt.Errorf("error in rendering password reset template: %w", err)
	}

	emailData := domain.EmailData{
		To:       payload.UserEmail,
		Subject:  "Recuperação de Senha - Store Manager",
		BodyHTML: htmlBody,
	}

	return p.sender.Send(ctx, emailData)
}

func (p *EmailProcessor) HandleEmailChangeTask(ctx context.Context, taskPayload []byte) error {
	var payload jobs.EmailChangeConfirmationJobData
	if err := json.Unmarshal(taskPayload, &payload); err != nil {
		return fmt.Errorf("unmarshal email verification payload: %w", err)
	}

	templateData := map[string]interface{}{
		"UserName":         payload.UserName,
		"ConfirmationLink": payload.ConfirmationLink,
	}

	htmlBody, err := p.templateManager.Render("email_verification.html", templateData)
	if err != nil {
		return fmt.Errorf("error in rendering email verification template: %w", err)
	}

	emailData := domain.EmailData{
		To:       payload.ToEmail,
		Subject:  "Confirmação de Alteração de Email - Store Manager",
		BodyHTML: htmlBody,
	}

	return p.sender.Send(ctx, emailData)
}

func (p *EmailProcessor) HandleAccountVerification(ctx context.Context, taskPayload []byte) error {
	var payload jobs.AccountVerificationJobData

	if err := json.Unmarshal(taskPayload, &payload); err != nil {
		return fmt.Errorf("unmarshal account verification payload: %w", err)
	}

	templateData := map[string]interface{}{
		"UserName":         payload.UserName,
		"VerificationLink": payload.VerificationLink,
	}

	htmlBody, err := p.templateManager.Render("account_verification.html", templateData)
	if err != nil {
		return fmt.Errorf("error in rendering account verification template: %w", err)
	}

	emailData := domain.EmailData{
		To:       payload.ToEmail,
		Subject:  "Verificação de Conta - Store Manager",
		BodyHTML: htmlBody,
	}

	return p.sender.Send(ctx, emailData)
}

func (p *EmailProcessor) HandlePromotionNotification(ctx context.Context, taskPayload []byte) error {
	var payload jobs.PromotionNotificationJobData

	if err := json.Unmarshal(taskPayload, &payload); err != nil {
		return fmt.Errorf("unmarshal promotion notification payload: %w", err)
	}

	templateData := map[string]interface{}{
		"UserName": payload.UserName,
		"NewRole":  payload.NewRole,
	}

	htmlBody, err := p.templateManager.Render("promotion_notification.html", templateData)
	if err != nil {
		return fmt.Errorf("error in rendering promotion notification template: %w", err)
	}

	emailData := domain.EmailData{
		To:       payload.ToEmail,
		Subject:  "Notificação de Promoção - Store Manager",
		BodyHTML: htmlBody,
	}

	return p.sender.Send(ctx, emailData)
}
