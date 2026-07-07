package validator

import govalidator "github.com/go-playground/validator/v10"

var validate = govalidator.New()

// FieldError is a single human-readable validation failure.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate validates a struct using its `validate` tags. It returns nil
// when the struct is valid, or a slice of field errors otherwise.
func Validate(s any) []FieldError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var out []FieldError
	for _, e := range err.(govalidator.ValidationErrors) {
		out = append(out, FieldError{
			Field:   e.Field(),
			Message: messageFor(e),
		})
	}
	return out
}

func messageFor(e govalidator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	case "email":
		return e.Field() + " must be a valid email"
	case "oneof":
		return e.Field() + " must be one of: " + e.Param()
	default:
		return e.Field() + " is invalid"
	}
}
