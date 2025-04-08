package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ValidationResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

var ErrValidator = map[string]string{}

func ErrValidationResponse(err error) []ValidationResponse {
	var validationResponses []ValidationResponse
	var fieldErrors validator.ValidationErrors

	if errors.As(err, &fieldErrors) {
		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				validationResponses = append(validationResponses, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is required", err.Field()),
				})
			case "email":
				validationResponses = append(validationResponses, ValidationResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is not a valid email address", err.Field()),
				})
			default:
				errValidator, ok := ErrValidator[err.Tag()]
				if ok {
					count := strings.Count(errValidator, "%s")
					if count == 1 {
						validationResponses = append(validationResponses, ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field()),
						})
					} else {
						validationResponses = append(validationResponses, ValidationResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field(), err.Param()),
						})
					}
				} else {
					validationResponses = append(validationResponses, ValidationResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("something is wrong on %s: %s", err.Field(), err.Tag()),
					})
				}
			}
		}
	}

	return validationResponses
}

func WrapError(err error) error {
	logrus.Error("error: %v", err)

	return err
}
