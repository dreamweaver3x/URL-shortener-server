package main

import (
	"SOKR/internal/app"
	"SOKR/internal/models"
	"SOKR/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
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
	go func() {
		for {
			application.CheckUrlStatus()
			time.Sleep(time.Minute * 10)
		}
	}()



	http.HandleFunc("/getshortstats", application.GetShortUrlStats)
	http.HandleFunc("/urlshortener", application.GetShortURL)
	http.HandleFunc("/", application.RedirectWithShortUrl)
	http.ListenAndServe(":8080", nil)
}
