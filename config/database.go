package config

import (
	"fmt"
	"log"
	"my_vocab/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DB_USER     = "myvocab"
	DB_PASSWORD = "myvocab123"
	DB_NAME     = "myvocab_db"
	DB_HOST     = "vocab-db.cbuoaypgqh0v.us-east-1.rds.amazonaws.com"
	DB_PORT     = "3306"
)

var DB *gorm.DB

func InitDb() {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	var err error
	DB, err = gorm.Open(mysql.Open(dbInfo), &gorm.Config{})

	if err != nil {
		log.Print("Connection Failed : ", err)
	}

	if err == nil {
		MigrateDB()
		log.Println("Connected!")
	}

}

func MigrateDB() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Vocab{})
	DB.AutoMigrate(&models.TypeVocab{})

}
