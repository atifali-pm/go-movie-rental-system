package db

import (
	"fmt"
	"os"

	"github.com/atifali-pm/go-movie-rental-system/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	godotenv.Load()

	DbHost := os.Getenv("MYSQL_HOST")
	DbName := os.Getenv("MYSQL_DBNAME")
	DbUsername := os.Getenv("MYSQL_USER")
	DbPassword := os.Getenv("MYSQL_PASSWORD")

	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbUsername, DbPassword, DbHost, DbName)
	dbConnection, err := gorm.Open(mysql.Open(connection), &gorm.Config{})

	if err != nil {
		panic("connection failed to the database")
	}

	DB = dbConnection
	fmt.Println("DB connected succesfully!")

	AutoMigrate(dbConnection)
}

func AutoMigrate(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&models.Category{},
		&models.Film{},
		&models.Inventory{},
		&models.Language{},
		&models.Actor{},
		&models.FilmActor{},
		&models.FilmCategory{},
		&models.Country{},
		&models.City{},
		&models.Address{},
		&models.Customer{},
		&models.Rental{},
		&models.Staff{},
		&models.Store{},
		&models.Payment{},
	)
}
