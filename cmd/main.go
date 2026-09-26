package main

import (
	"log"

	"exam/app"
	"exam/config"
	"exam/database"
)

// @title Exam Product API
// @version 1.0
// @description API for managing products (create, update)
// @host localhost:3000
// @BasePath /
func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg.DatabaseURL)
	database.Migrate(db)

	app := app.NewApp(db)

	log.Fatal(app.Listen(":" + cfg.Port))
}
