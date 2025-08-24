package handlers

import (
	"context"
	"log"
	"sica/internal/dtos"
	"sica/internal/models"
	"sica/internal/repositories"
	"sica/pkg/cloudinary"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"

	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {

	r := repositories.NewProductRepository()
	products, err := r.GetAll()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error getting the products",
		})
		return
	}

	dtoProducts := dtos.ToProductsDetailsResponse(products)

	c.JSON(200, gin.H{
		"products": dtoProducts,
	})
}

func GetProduct(c *gin.Context) {

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid id format",
		})
		return
	}
	id := uint(idUint64)

	r := repositories.NewProductRepository()
	product, err := r.Get(id)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error getting product",
		})
		return
	}

	dtoProduct := dtos.ToProductDetailsResponse(&product)

	c.JSON(200, gin.H{
		"product": dtoProduct,
	})
}

func CreateProduct(c *gin.Context) {

	var product models.Product
	err := c.BindJSON(&product)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid json format",
		})
		return
	}

	if product.Name == "" {
		c.JSON(400, gin.H{
			"error": "product name can't be empty",
		})
		return
	}

	r := repositories.NewProductRepository()
	newProduct, err := r.Create(&product)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error creating product",
		})
		return
	}

	dtoProduct := dtos.ToProductDetailsResponse(newProduct)

	c.JSON(201, gin.H{
		"message": "product created",
		"product": dtoProduct,
	})
}

func UpdateProduct(c *gin.Context) {

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid id format",
		})
		return
	}
	id := uint(idUint64)

	var product models.Product
	err = c.BindJSON(&product)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid json format",
		})
		return
	}

	if product.Name == "" {
		c.JSON(400, gin.H{
			"error": "product name can't be empty",
		})
		return
	}

	newValues := map[string]interface{}{
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"category_id": product.CategoryID,
		"available":   product.Available,
		"visible":     product.Visible,
		"image":       product.Image,
	}

	r := repositories.NewProductRepository()
	newProduct, err := r.Update(id, newValues)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error updating product",
		})
		return
	}

	dtoProduct := dtos.ToProductDetailsResponse(&newProduct)

	c.JSON(200, gin.H{
		"message": "product updated",
		"product": dtoProduct,
	})
}

func DeleteProduct(c *gin.Context) {

	idStr := c.Param("id")
	idUint64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid id format",
		})
		return
	}
	id := uint(idUint64)

	r := repositories.NewProductRepository()
	err = r.Delete(id)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error deleting product",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "product deleted",
	})
}

func DeleteImage(c *gin.Context) {
	publicID := c.Param("id")
	if publicID == "" {
		c.JSON(400, gin.H{"error": "public id required"})
		return
	}

	cld, err := cloudinary.GetCloudinary()
	if err != nil {
		c.JSON(500, gin.H{"error": "cloudinary not initialized"})
		return
	}

	ctx := context.Background()
	_, err = cld.C.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "error deleting image from cloudinary"})
		return
	}

	c.JSON(200, gin.H{
		"message": "image deleted",
	})
}
