package handlers

import (
	"log"
	"sica/internal/dtos"
	"sica/internal/models"
	"sica/internal/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllCategories(c *gin.Context) {
	r := repositories.NewCategoryRepository()
	categories, err := r.GetAll()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error getting the categories",
		})
		return
	}

	dtoCategories := dtos.ToCategoriesSimpleResponse(categories)

	c.JSON(200, gin.H{
		"categories": dtoCategories,
	})
}

func GetAllCP(c *gin.Context) {
	r := repositories.NewCategoryRepository()
	categories, err := r.GetAllCP()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error getting the categories",
		})
		return
	}

	dtoCategories := dtos.ToCategoriesDetailsResponse(categories)

	c.JSON(200, gin.H{
		"categories": dtoCategories,
	})
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	err := c.BindJSON(&category)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid json format",
		})
		return
	}

	if category.Name == "" {
		c.JSON(400, gin.H{
		"error": "category name can't be empty",
	})
	return
	}

    if category.Order != 0 || category.InternalUpdate  || category.Products != nil {

		c.JSON(400, gin.H{
			"error": "too many params",
		})
		return
    }

	r := repositories.NewCategoryRepository()
	newCategory, err := r.Create(&category)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error creating category",
		})
		return
	}

	dtoCategory := dtos.ToCategorySimpleResponse(newCategory)

	c.JSON(201, gin.H{
		"message":  "category created",
		"category": dtoCategory,
	})
}

func UpdateCategory(c *gin.Context) {
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

	var category models.Category
	err = c.BindJSON(&category)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": "invalid json format",
		})
		return
	}

	newValues := make(map[string]interface{})

	if category.Name != "" {
		newValues["name"] = category.Name
	}
	if category.Order != 0 {
		newValues["order"] = category.Order
	}

	if len(newValues) == 0 {
		c.JSON(200, gin.H{
			"message": "nothing to update",
		})
		return
	}

	r := repositories.NewCategoryRepository()
	updatedCategory, err := r.Update(id, newValues)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error updating category",
		})
		return
	}

	dtoCategory := dtos.ToCategorySimpleResponse(&updatedCategory)

	c.JSON(200, gin.H{
		"message":  "category updated",
		"category": dtoCategory,
	})
}

func DeleteCategory(c *gin.Context) {
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

	r := repositories.NewCategoryRepository()
	err = r.Delete(id)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "error deleting category",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "category deleted",
	})
}
