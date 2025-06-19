package ports

type TemplateManager interface {
	Render(templateName string, data interface{}) (string, error)
}
