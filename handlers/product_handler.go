package handlers

import (
	"context"
	"net/http"
	"products/config"
	"products/models"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateProduct(c echo.Context) error {
	product := new(models.Product)
	if err := c.Bind(product); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request body",
			"msg":   err.Error(),
		})
	}

	collection := config.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "failed to insert product",
			"msg":   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "product created successfully",
		"data":    result,
	})
}

func GetProduct(c echo.Context) error {
	id := c.Param("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid product ID format",
		})
	}

	collection := config.GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product models.Product
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "product not found",
			"msg":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, product)
}
