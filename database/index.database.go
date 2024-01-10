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

	dsnMysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "", "127.0.0.1", "3306", "kulina_test")

	DB, ConnectionError = gorm.Open(mysql.Open(dsnMysql), &gorm.Config{})

	if ConnectionError != nil {
		panic("Cant connect to database")
	}

	DB.AutoMigrate(
		&model.User{}, 
		&model.Address{}, 
		&model.Supplier{},
		&model.Product{},
		&model.Store{},
		&model.StoreSellingArea{},
		&model.StoreProduct{},
	)

	if err := SeedProducts(DB); err != nil {
		fmt.Println("Error seeding products:", err)
		return
	}

	log.Println("Connect database")

}

func SeedProducts(db *gorm.DB) error {
	
	var count int64
	db.Model(&model.Product{}).Count(&count)

	if count > 0 {
		return nil
	}

	products := []model.Product{
		{Name: "Product A", Price: 50000},
		{Name: "Product B", Price: 40000},
		{Name: "Product C", Price: 30000},
		{Name: "Product D", Price: 50000},
		{Name: "Product E", Price: 40000},
		{Name: "Product F", Price: 30000},
	}

	result := db.Create(&products)
	if result.Error != nil {
		return result.Error
	}

	return nil
}