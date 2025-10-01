package response

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Translate validation error agar user-friendly
func TranslateValidationError(err error) string {
	var sb strings.Builder

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := e.Field()
			switch e.Tag() {
			case "required":
				sb.WriteString(field + " is required. ")
			case "email":
				sb.WriteString("Invalid email format. ")
			case "min":
				sb.WriteString(field + " must be at least " + e.Param() + " characters. ")
			case "max":
				sb.WriteString(field + " must be at most " + e.Param() + " characters. ")
			default:
				sb.WriteString(field + " is invalid. ")
			}
		}
		return strings.TrimSpace(sb.String())
	}

	return err.Error()
}
