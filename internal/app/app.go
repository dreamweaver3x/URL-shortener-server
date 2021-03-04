package app

import (
	"SOKR/internal/models"
	"SOKR/internal/repository"
	"log"
	"net/http"
	_ "net/url"
	_ "strings"
)

type Application struct {
	repo *repository.LinksRepository
}

func NewApplication(repo *repository.LinksRepository) *Application {
	return &Application{repo: repo}
}

func (a *Application) GetShortURL(w http.ResponseWriter, r *http.Request) {
	u := &models.Link{FullUrl: r.URL.Query().Get("link")}
	hasLong := a.repo.CheckForElemLong(u)
	if hasLong {
		log.Println("uzhe est' elem")
		return
	} else {
		log.Println("elema net")
	}
	shortUrl, err := a.repo.Create(u)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write([]byte("http://localhost:8080/" + shortUrl.ShortUrl))
	w.WriteHeader(http.StatusOK)
}

func (a *Application) RedirectWithShortUrl(w http.ResponseWriter, r *http.Request) {
	u := &models.Link{}
	linkWithSlash := r.URL.String()
	u.ShortUrl = linkWithSlash[1:]
	hasShort := a.repo.CheckForElemShort(u)
	if !hasShort {
		println("don't have this url")
		return
	}
	u, err := a.repo.GetLongUrl(u)
	if err != nil {
		log.Println(err)
		return
	}
	println(u.FullUrl)
	http.Redirect(w, r, "https://" + u.FullUrl, http.StatusMovedPermanently)
}
