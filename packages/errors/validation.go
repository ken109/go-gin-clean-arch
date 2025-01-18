package errors

import (
	"encoding/json"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"

	"go-gin-clean-arch/packages/validation"
)

type Validation struct {
	errors map[string][]string
}

func NewValidation() *Error {
	return &Error{
		kind:       KindValidation,
		validation: &Validation{errors: map[string][]string{}},
	}
}

func (verr *Validation) MarshalJSON() ([]byte, error) {
	return json.Marshal(verr.errors)
}

func (verr *Validation) Add(fieldName string, message string) {
	if _, ok := verr.errors[fieldName]; !ok {
		verr.errors[fieldName] = []string{message}
	} else {
		verr.errors[fieldName] = append(verr.errors[fieldName], message)
	}
}

func (verr *Validation) Error() string {
	b, _ := json.Marshal(verr)
	return string(b)
}

var getFieldPathRegex = regexp.MustCompile(`^[^.]+\.([^.]+\..+)`)

func (verr *Validation) Validate(request interface{}) (invalid bool) {
	err := validation.Validate().Struct(request)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			for _, f := range err.(validator.ValidationErrors) {
				fieldPath := getFieldPathRegex.FindStringSubmatch(f.StructNamespace())
				if len(fieldPath) > 0 {
					verr.Add(strcase.ToSnakeWithIgnore(fieldPath[1], "."), f.Translate(validation.Translator()))
				} else {
					verr.Add(strcase.ToSnakeWithIgnore(f.StructField(), "."), f.Translate(validation.Translator()))
				}
			}
		}
	}
	return verr.IsInvalid()
}

func (verr *Validation) IsInvalid() bool {
	return len(verr.errors) > 0
}
