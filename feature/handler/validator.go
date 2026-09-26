package producthandler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func FormatValidationErrors(err error) []fiber.Map {
	var errors []fiber.Map

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range validationErrors {
			errors = append(errors, fiber.Map{
				"field": fe.Field(),
				"tag":   fe.Tag(),
				"value": fe.Param(),
			})
		}
	}

	return errors
}
