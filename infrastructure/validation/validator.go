package web_errors

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslation "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate = validator.New()
	transl   ut.Translator
)

func init() {
	en := en.New()
	unt := ut.New(en, en)
	transl, _ = unt.GetTranslator("en")
	entranslation.RegisterDefaultTranslations(Validate, transl)
}

func TranslateError(err error) *RestErr {
	var jsonErr *json.UnmarshalTypeError
	var validationErrs validator.ValidationErrors

	if errors.As(err, &jsonErr) {
		return NewBadRequestError("Invalid field type")
	} else if errors.As(err, &validationErrs) {
		var causes []Causes

		for _, e := range validationErrs {
			cause := Causes{
				Message: e.Translate(transl),
				Field:   e.Field(),
			}

			causes = append(causes, cause)
		}

		return NewBadRequestValidationError("Some fields are invalid", causes)
	} else {
		return NewBadRequestError("Error trying to convert fields")
	}
}
