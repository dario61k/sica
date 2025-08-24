package repositories

import (
	"sica/internal/database"
	"sica/internal/models"
	"sync"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

var onceCR sync.Once
var singletonCategoryRepository *categoryRepository

func NewCategoryRepository() *categoryRepository {
	onceCR.Do(func() {
		singletonCategoryRepository = &categoryRepository{db: database.GetDB()}
	})
	return singletonCategoryRepository
}

func (r *categoryRepository) Create(category *models.Category) (*models.Category, error) {
	if result := r.db.Model(&models.Category{}).Create(category); result.Error != nil {
		return nil, result.Error
	}
	if err := r.db.First(category, category.ID).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) GetAll() ([]models.Category, error) {
	var categories []models.Category
	result := r.db.Model(&models.Category{}).Order(`"order"`).Find(&categories)
	if result.Error != nil {
		return []models.Category{}, result.Error
	}
	return categories, nil
}

func (r *categoryRepository) GetAllCP() ([]models.Category, error) {

	var categories []models.Category
	result := r.db.Model(&models.Category{}).
	Joins("JOIN products ON products.category_id = categories.id").
	Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "image", "name", "price", "available", "description", "category_id").
				Where("visible = ?", true)
	}).
	Group("categories.id").
	Order(`"order"`).
	Find(&categories)


	if result.Error != nil {
		return []models.Category{}, result.Error
	}
	return categories, nil
}

func (r *categoryRepository) Update(id uint, newValues map[string]interface{}) (models.Category, error) {

	result := r.db.Model(&models.Category{}).Where("id = ?", id).Updates(newValues)
	if result.Error != nil {
		return models.Category{}, result.Error
	}
	var category models.Category
	result = r.db.Model(&models.Category{}).First(&category, id)
	if result.Error != nil {
		return models.Category{}, result.Error
	}
	return category, nil
}

func (r *categoryRepository) Delete(id uint) error {
	if result := r.db.Model(&models.Category{}).Delete(&models.Category{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}
