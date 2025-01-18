package util

import (
	"context"

	"go-gin-clean-arch/config"
	"go-gin-clean-arch/packages/errors"
)

func getValidationError(ctx context.Context) *errors.Error {
	validationError, _ := ctx.Value(config.ErrorKey).(*errors.Error)
	return validationError
}

func Validate(ctx context.Context, request interface{}) (invalid bool) {
	validationError := getValidationError(ctx)

	return validationError.Validation().Validate(request)
}

func InvalidField(ctx context.Context, fieldName string, reason string) {
	validationError := getValidationError(ctx)

	validationError.Validation().Add(fieldName, reason)
}

func IsInvalid(ctx context.Context) bool {
	validationError := getValidationError(ctx)

	return validationError.Validation().IsInvalid()
}

func ValidationError(ctx context.Context) error {
	validationError := getValidationError(ctx)

	return validationError
}
