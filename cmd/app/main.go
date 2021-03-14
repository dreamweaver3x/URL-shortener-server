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
	"github.com/gorilla/mux"
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
			application.CheckUrlStatusNew()
			time.Sleep(time.Minute * 10)
		}
	}()

r := mux.NewRouter()

	r.HandleFunc("/shortstats", application.GetShortUrlStats)
	r.HandleFunc("/urlshortener", application.GetShortURL).Methods("POST")
	r.HandleFunc("/", application.RedirectWithShortUrl)
	http.ListenAndServe(":8080", r)
}
