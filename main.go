package main

import (
	_ "REST/docs"
	"REST/routes"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/books", routes.GetBooks)
		api.GET("/books/:id", routes.GetBook)
		api.POST("/books", routes.AddBook)
		api.PUT("/books/:id", routes.UpdateBook)
		api.DELETE("/books/:id", routes.DeleteBook)
	}

	// Register Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
