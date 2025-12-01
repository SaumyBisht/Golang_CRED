// main.go
package main

import (
	"products/config"
	"products/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Database connection
	config.ConnectDB()

	// Routes
	e.POST("/products", handlers.CreateProduct)
	e.GET("/products/:id", handlers.GetProduct)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
