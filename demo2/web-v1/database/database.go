package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect2DB() *gorm.DB {
	dsn := "root:chen0309@tcp(localhost:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
	fmt.Printf("connect to Database: %s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Failed when connected to database, error= %v", err)
		return nil
	}
	return db
}
