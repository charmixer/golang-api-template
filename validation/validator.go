package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

const locale = "en"

var (
	Validate    *validator.Validate
	Translation ut.Translator
)

func init() {
	Validate = validator.New()

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	enTranslator := en.New()
	uni := ut.New(enTranslator, enTranslator)

	Translation, _ = uni.GetTranslator(locale)
	en_translations.RegisterDefaultTranslations(Validate, Translation)
}
