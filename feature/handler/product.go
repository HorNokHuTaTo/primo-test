package producthandler

import (
	productmodel "exam/feature/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// CreateProduct godoc
// @Summary Create a new product
// @Description Creates a product with name, price, and optional description/sale_price
// @Tags products
// @Accept json
// @Produce json
// @Param request body productmodel.CreateProductRequest true "Product to create"
// @Success 201 {object} productmodel.CreateSuccessResponse "Successful creation"
// @Failure 400 {object} productmodel.ErrorResponse "BODY_PARSER_ERROR, VALIDATOR_ERROR"
// @Failure 500 {object} productmodel.ErrorResponse "INTERNAL_SERVER_ERROR"
// @Router /product [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	ctx := c.UserContext()
	req := productmodel.CreateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error_code": "BODY_PARSER_ERROR",
			"errors":     err.Error(),
		})
	}

	if err := validator.New().Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error_code": "VALIDATOR_ERROR",
			"errors":     FormatValidationErrors(err),
		})
	}

	res, err := h.ProductService.CreateProduct(ctx, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"successful": false,
			"error_code": "INTERNAL_SERVER_ERROR",
			"errors":     err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"successful": true,
		"data":       res,
	})
}

// UpdateProduct godoc
// @Summary Partially update a product
// @Description Partial update — only send fields you want to change.
// @Description "name" and "price" cannot be cleared, only updated or omitted.
// @Description "description" and "sale_price" can be omitted (no change), set to null (clear the value), or set to a new value.
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Param request body productmodel.UpdateProductRequest true "Fields to update"
// @Success 200 {object} productmodel.UpdateSuccessResponse "Successful update"
// @Failure 400 {object} productmodel.ErrorResponse "MISSING_PRODUCT_ID, INVALID_PRODUCT_ID, BODY_PARSER_ERROR, or VALIDATOR_ERROR"
// @Failure 500 {object} productmodel.ErrorResponse "INTERNAL_SERVER_ERROR"
// @Router /product/{id} [patch]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	ctx := c.UserContext()
	productID := c.Params("id")
	validate := validator.New()
	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error_code": "MISSING_PRODUCT_ID",
			"errors":     "Product ID is required",
		})
	} else {
		if err := validate.Var(productID, "uuid"); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"successful": false,
				"error_code": "INVALID_PRODUCT_ID",
				"errors":     "Product ID must be a valid UUID",
			})
		}
	}

	req := productmodel.UpdateProductRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error_code": "BODY_PARSER_ERROR",
			"errors":     err.Error(),
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"successful": false,
			"error_code": "VALIDATOR_ERROR",
			"errors":     FormatValidationErrors(err),
		})
	}

	if err := h.ProductService.UpdateProduct(ctx, productID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"successful": false,
			"error_code": "INTERNAL_SERVER_ERROR",
			"errors":     err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"successful": true,
	})
}
