package repositories

import (
	"sica/internal/database"
	"sica/internal/models"
	"sync"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

var oncePR sync.Once
var singletonProductRepository *productRepository

func NewProductRepository() *productRepository {
	oncePR.Do(func() {
		singletonProductRepository = &productRepository{db: database.GetDB()}
	})
	return singletonProductRepository
}

func (r *productRepository) Create(product *models.Product) (*models.Product, error) {
	if result := r.db.Model(&models.Product{}).Create(product); result.Error != nil {
		return nil, result.Error
	}
	if err := r.db.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).First(product, product.ID).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	result := r.db.Model(&models.Product{}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).Order("id DESC").Find(&products)
	if result.Error != nil {
		return []models.Product{}, result.Error
	}
	return products, nil
}

func (r *productRepository) Get(id uint) (models.Product, error) {
	var products models.Product
	result := r.db.Model(&models.Product{}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).First(&products, id)
	if result.Error != nil {
		return models.Product{}, result.Error
	}
	return products, nil
}

func (r *productRepository) Update(id uint, newValues map[string]interface{}) (models.Product, error) {
	result := r.db.Model(&models.Product{}).Where("id = ?", id).Updates(newValues)
	if result.Error != nil {
		return models.Product{}, result.Error
	}
	var product models.Product
	result = r.db.Model(&models.Product{}).Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name")
	}).First(&product, id)
	if result.Error != nil {
		return models.Product{}, result.Error
	}
	return product, nil
}

func (r *productRepository) Delete(id uint) error {
	if result := r.db.Model(&models.Product{}).Delete(&models.Product{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}
