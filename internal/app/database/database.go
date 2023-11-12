package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"github.com/nukkua/ra-chi/internal/app/models"
)


func SetupDatabase () * gorm.DB {
	
	
	requiredEnvVars := []string{"DB_USER","DB_PASSWORD", "DB_HOST", "DB_NAME"}

	for _, envVar := range requiredEnvVars{
		if os.Getenv(envVar) == ""{
			log.Fatalf("This env variable is required but not defined: %s", envVar);
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_HOST"),
	os.Getenv("DB_NAME"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to database");
	}

	db.AutoMigrate(&models.User{})

	return db
}
