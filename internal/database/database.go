package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Groups struct {
	gorm.Model
	Batch   string
	Names   string
	IsBase  bool
	Project string
	Group1  string
	Group2  string
	Group3  string
	Group4  string
	Group5  string
	Group6  string
	Group7  string
}

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open("groups.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&Groups{})
	fmt.Println("Connection Opened to Database and Migrated")
	return DB
}
