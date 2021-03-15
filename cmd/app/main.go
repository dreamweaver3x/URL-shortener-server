package main

import (
	"SOKR/config"
	"SOKR/internal/app"
	"SOKR/internal/models"
	"SOKR/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	conf := config.Load()
	db, err := gorm.Open(postgres.Open(conf.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = models.InitModels(db)
	if err != nil {
		log.Fatal("too bad")
	}
	repo := repository.NewLinksRepository(db)
	application := app.NewApplication(repo)
	application.Start(conf.Port)
}
