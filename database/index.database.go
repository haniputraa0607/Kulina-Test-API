package database

import (
	"fmt"
	"log"
	"rest-api/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	var ConnectionError error

	dsnMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "12345678", "127.0.0.1", "3306", "kulina_test")

	DB, ConnectionError = gorm.Open(mysql.Open(dsnMysql), &gorm.Config{})

	if ConnectionError != nil {
		panic("Cant connect to database")
	}

	DB.AutoMigrate(
		&model.User{}, 
		&model.Address{}, 
		&model.Supplier{},
		&model.Store{},
		&model.StoreSellingArea{},
		&model.Product{},
	)

	log.Println("Connect database")

}