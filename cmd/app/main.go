package main

import (
	"SOKR/config"
	"SOKR/internal/app"
	"SOKR/internal/models"
	"SOKR/internal/repository"
	"flag"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dev := flag.Bool("dev",
		false,
		"enable reading config from .env file instead of system env vars",
	)
	flag.Parse()

	if *dev {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gorm.Open(postgres.Open(conf.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = models.InitModels(db)
	if err != nil {
		log.Fatal("too bad")
	}
	repo := repository.NewLinksRepository(db)
	application := app.NewApplication(repo)
	application.Start(conf.ListenAddress())
}
