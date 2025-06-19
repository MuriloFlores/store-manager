package templates

import (
	"bytes"
	"fmt"
	"github.com/muriloFlores/StoreManager/internal/core/ports"
	"html/template"
)

type HTMLTemplateManager struct {
	templates *template.Template
}

func NewHTMLTemplateManager(templatesDir string) (ports.TemplateManager, error) {
	templates, err := template.ParseGlob(fmt.Sprintf("%s/*.html", templatesDir))
	if err != nil {
		return nil, fmt.Errorf("template parsing error: %s", err)
	}

	return &HTMLTemplateManager{templates: templates}, nil
}

func (m *HTMLTemplateManager) Render(templateName string, data interface{}) (string, error) {
	var body bytes.Buffer

	if err := m.templates.ExecuteTemplate(&body, templateName, data); err != nil {
		return "", fmt.Errorf("template rendering error: %s", err)
	}

	return body.String(), nil
}
