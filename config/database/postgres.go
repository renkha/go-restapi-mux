package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

func DbInit() *gorm.DB {
	log.Info("Startting Database Connection")

	dbURI := "host=localhost user=postgres password=postgres dbname=todo_list port=5432 sslmode=disable"

	db, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Panic("Failed to connect database with error" + err.Error())
	}

	db.LogMode(true)

	return db
}
