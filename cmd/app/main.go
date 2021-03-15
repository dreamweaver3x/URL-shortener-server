package main

import (
	"SOKR/internal/app"
	"SOKR/internal/models"
	"SOKR/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {

	dsn := "host=localhost user=db_user password=pwd123 dbname=urlcutter port=54320 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = models.InitModels(db)
	if err != nil {
		log.Fatal("too bad")
	}
	repo := repository.NewLinksRepository(db)
	application := app.NewApplication(repo)
	application.Start()
}
