package validator

import (
	"errors"
	"path"

	applog "github.com/GBA-BI/tes-api/pkg/log"
	"github.com/go-playground/validator/v10"

	apperrors "github.com/GBA-BI/tes-api/pkg/errors"
)

var gValidator = validator.New()

func init() {
	RegisterValidators()
}

// RegisterValidators registers k8s validators
func RegisterValidators() {
	validates := []struct {
		tag string
		fn  validator.Func
	}{
		{tag: "abspath", fn: validateAbsPath},
	}

	for _, v := range validates {
		if err := gValidator.RegisterValidation(v.tag, v.fn); err != nil {
			applog.Fatalw("register custom validation", "name", v.tag, "err", err)
		}
	}
}

// Validate can validate struct field with validate tag
func Validate(s interface{}) error {
	err := gValidator.Struct(s)
	if err != nil {
		applog.Errorw("validation error", "err", err)
		if validationErrors := make(validator.ValidationErrors, 0); errors.As(err, &validationErrors) {
			fields := make([]string, 0, len(validationErrors))
			for _, validationErr := range validationErrors {
				fields = append(fields, validationErr.Field())
			}
			return apperrors.NewInvalidError(fields...)
		}
		return apperrors.NewInvalidError()
	}
	return nil
}

func validateAbsPath(fl validator.FieldLevel) bool {
	return path.IsAbs(fl.Field().String())
}
