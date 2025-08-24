package repositories

import (
	"sica/internal/database"
	"sica/internal/models"
	"sync"

	"gorm.io/gorm"
)

type authRepository struct {
    db *gorm.DB
}

var onceAuth sync.Once
var singletonAuthRepository *authRepository

func NewAuthRepository() *authRepository {
    onceAuth.Do(func() {
        singletonAuthRepository = &authRepository{db: database.GetDB()}
    })
    return singletonAuthRepository
}

func (r *authRepository) Get(id uint) (models.Auth, error) {
    var auth models.Auth
    result := r.db.Model(&models.Auth{}).First(&auth, id)
    if result.Error != nil {
        return models.Auth{}, result.Error
    }
    return auth, nil
}

func (r *authRepository) Create(auth *models.Auth) (*models.Auth, error) {
    if result := r.db.Model(&models.Auth{}).Create(auth); result.Error != nil {
        return nil, result.Error
    }
    return auth, nil
}

func (r *authRepository) Update(id uint, newValues map[string]interface{}) (models.Auth, error) {
    result := r.db.Model(&models.Auth{}).Where("id = ?", id).Updates(newValues)
    if result.Error != nil {
        return models.Auth{}, result.Error
    }
    var auth models.Auth
    result = r.db.Model(&models.Auth{}).First(&auth, id)
    if result.Error != nil {
        return models.Auth{}, result.Error
    }
    return auth, nil
}

func (r *authRepository) Delete(id uint) error {
    if result := r.db.Model(&models.Auth{}).Delete(&models.Auth{}, id); result.Error != nil {
        return result.Error
    }
    return nil
}