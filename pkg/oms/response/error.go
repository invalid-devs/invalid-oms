package response

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
	"reflect"
	"strings"
)

type Error struct {
	Message     string              `json:"message"`
	FieldErrors []map[string]string `json:"field_errors"`
}

func Validate(r *http.Request, s any) *Error {
	body, _ := io.ReadAll(r.Body)

	if err := json.Unmarshal(body, &s); err != nil {
		msg := Error{
			Message: err.Error(),
		}

		return &msg
	}

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	if err := validate.Struct(s); err != nil {
		var validationErrors []map[string]string
		for _, fieldError := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, map[string]string{
				fieldError.Field(): fieldError.Tag(),
			})
		}
		msg := Error{
			Message:     "Validation errors!",
			FieldErrors: validationErrors,
		}
		return &msg
	}

	return nil
}
