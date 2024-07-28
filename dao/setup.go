package dao

import (
	"RuleEngineAST/models"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase(dbName string) {

	database, err := gorm.Open(sqlite.Open(dbName+"test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.Rule{})
	if err != nil {
		return
	}

	DB = database
}
