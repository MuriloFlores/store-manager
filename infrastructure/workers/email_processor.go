package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/domain"
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
	var payload domain.PasswordChangeJobData
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
	var payload domain.EmailChangeConfirmationJobData
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
