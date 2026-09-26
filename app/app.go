package app

import (
	producthandler "exam/feature/handler"
	productrepository "exam/feature/repository"
	productservice "exam/feature/service"

	_ "exam/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/gofiber/swagger"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	SetupRoutes(app, db)

	return app
}

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	productRepository := productrepository.NewProductRepository(db)
	productService := productservice.NewProductService(productRepository)
	productHandler := producthandler.NewProductHandler(productService)
	app.Get("/api-docs/*", fiberSwagger.HandlerDefault) // Swagger UI route
	app.Group("/product").
		Post("/", productHandler.CreateProduct).
		Patch("/:id", productHandler.UpdateProduct)
}
