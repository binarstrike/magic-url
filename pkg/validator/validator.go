package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// TODO: kembalikan data kesalahan setiap field dalam bentuk json

type Validator struct {
	Validate *validator.Validate
}

func NewValidator() *Validator {
	v := new(Validator)
	v.Validate = validator.New(validator.WithRequiredStructEnabled())
	return v
}

func (v Validator) ValidateStruct(s any) error {
	err := v.Validate.Struct(s)
	if err != nil {
		return parseValidationError(err)
	}

	return nil
}

type ValidationError struct {
	FailedField string
	Tag         string
	Value       any
}

func (ve *ValidationError) Error() string {
	if ve == nil {
		return ""
	}

	// TODO: pertimbangkan lagi teks kesalahan agar tidak menampilkan kredensial pengguna ketika terjadi kesalahan validasi
	return fmt.Sprintf("field `%s` with value '%v' needs to implement '%s'", ve.FailedField, ve.Value, ve.Tag)
}

func parseValidationError(e error) *ValidationError {
	if err, ok := e.(validator.ValidationErrors); ok && len(err) > 0 {
		return &ValidationError{
			FailedField: err[0].Field(),
			Tag:         err[0].Tag(),
			Value:       err[0].Value(),
		}
	}

	return nil
}
