package internal

import (
	"sica/internal/handlers"
	"sica/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine{

	r := gin.Default()

	r.Use(middleware.SetCors())

	r.HEAD("/health", func(c *gin.Context) { // health check
		c.Status(200)
	})

	api := r.Group("/api")
	p := api.Group("/product")
	c := api.Group("/category") 
	p.Use(middleware.AuthMiddleware())
	c.Use(middleware.AuthMiddleware())
	
	api.GET("/auth", middleware.AuthMiddleware(),func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message" : "auth",
		})
	})

	api.GET("/get-all", handlers.GetAllCP)

	p.GET("/", handlers.GetAllProducts)
	p.GET("/:id", handlers.GetProduct)
	p.POST("/", handlers.CreateProduct)
	p.PUT("/:id", handlers.UpdateProduct)
	p.DELETE("/:id", handlers.DeleteProduct)
	p.DELETE("/image/:id", handlers.DeleteImage)

	c.GET("/", handlers.GetAllCategories) 
	c.POST("/", handlers.CreateCategory)
	c.PUT("/:id", handlers.UpdateCategory)
	c.DELETE("/:id", handlers.DeleteCategory)

	
	api.POST("/login", handlers.Login)
	api.POST("/refresh-token", handlers.Refresh)
	api.PUT("/change-password", middleware.AuthMiddleware(), handlers.ChangePassword)



	return r

}