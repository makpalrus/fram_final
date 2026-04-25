package config

import (
	"final/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=Kz123456 dbname=fram port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = db.AutoMigrate(&models.User{}, &models.Course{}, &models.Enrollment{})
	if err != nil {
		log.Println("Migration failed:", err)
	}

	DB = db
	log.Println("Database connected and migrated!")
}
