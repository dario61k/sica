package database

import (
	"log"
	"os"
	"sica/internal/models"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		dsn := os.Getenv("DB_URL")
		log.Println(dsn)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Database conection failed", err)
		} 
	})
	return db
}

func AutoMigrate() {
	db.AutoMigrate(models.Category{}, models.Product{}, models.Auth{})
}
