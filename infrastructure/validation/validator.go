package validation

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/muriloFlores/StoreManager/infrastructure/web/web_errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translation "github.com/go-playground/validator/v10/translations/en"
)

var (
	// Validate é a instância singleton do nosso validador.
	Validate = validator.New()
	transl   ut.Translator
)

// init configura o validador e as traduções na inicialização do pacote.
func init() {
	en := en.New()
	unt := ut.New(en, en)
	transl, _ = unt.GetTranslator("en")
	en_translation.RegisterDefaultTranslations(Validate, transl)
}

// TranslateError traduz erros da biblioteca 'validator' para o nosso RestErr.
func TranslateError(err error) *web_errors.RestErr {
	var jsonErr *json.UnmarshalTypeError
	var validationErrs validator.ValidationErrors

	if errors.As(err, &jsonErr) {
		return web_errors.NewBadRequestError("Tipo de campo inválido, verifique os dados enviados.")
	} else if errors.As(err, &validationErrs) {
		causes := []web_errors.Causes{}

		for _, e := range validationErrs {
			cause := web_errors.Causes{
				Message: e.Translate(transl),
				Field:   e.Field(),
			}
			causes = append(causes, cause)
		}
		return web_errors.NewBadRequestValidationError("Alguns campos são inválidos", causes)
	} else {
		return web_errors.NewBadRequestError("Erro ao tentar converter os campos da requisição.")
	}
}
