package main

import (
	"SOKR/internal/app"
	"SOKR/internal/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {

	dsn := "host=localhost user=db_user password=pwd123 dbname=urlcutter port=54320 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err1 := data.InitModels(db)
	if err1 != nil {
		log.Fatal("too bad")
	}
	//data.ShortUrl(db)
	//data.CheckForElemLong(db, "http://google.com")
	//
	//data.Create(db)
	//
	http.HandleFunc("/urlshortener", app.GetShortURL)
	http.HandleFunc("/", app.RedirectWithShortUrl)
	http.ListenAndServe(":8080", nil)
}
