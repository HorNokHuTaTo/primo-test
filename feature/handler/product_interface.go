package producthandler

import (
	productservice "exam/feature/service"

	"github.com/gofiber/fiber/v2"
)

type IProductHandler interface {
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
}

type ProductHandler struct {
	ProductService productservice.IProductService
}

func NewProductHandler(productService productservice.IProductService) IProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}
