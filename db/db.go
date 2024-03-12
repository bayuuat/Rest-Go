package db

import (
	"RestAPI/model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := "127.0.0.1"
	port := "5432"
	user := "postgres"
	password := "password"
	dbname := "postgres"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrasi tabel
	DB.AutoMigrate(&model.Order{}, &model.Item{})
}
