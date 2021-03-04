package app

import (
	"SOKR/internal/data"
	_ "SOKR/internal/shorturl"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	_ "net/url"
	_ "strings"
)

const dsn = "host=localhost user=db_user password=pwd123 dbname=urlcutter port=54320 sslmode=disable"

func GetShortURL(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	link := r.URL.Query().Get("link")
	println(link)
	hasLong := data.CheckForElemLong(db, link)
	if hasLong {
		println("uzhe est' elem")
		return
	} else {
		println("elema net")
	}
	shortUrl := data.AddElem(db, link)
	fmt.Fprint(w, shortUrl)

}

func RedirectWithShortUrl(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	linkWithSlash := r.URL.String()
	println(linkWithSlash)
	link := linkWithSlash[1:]
	hasShort := data.CheckForElemShort(db, link)
	if !hasShort {
		println("don't have this url")
		return
	}
	fullLink := data.GetLongUrl(db, link)
	if fullLink != "" {
		http.Redirect(w, r, fullLink, http.StatusMovedPermanently)
	}
}
