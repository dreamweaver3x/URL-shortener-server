package app

import (
	"SOKR/internal/models"
	"SOKR/internal/repository"
	"encoding/json"
	"log"
	"net/http"
	_ "net/url"
	_ "strings"
	"sync"
)

type UrlSlice struct {
	sync.Mutex
	idSlice []uint
}

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
		w.Write([]byte("http://localhost:8080/" + u.ShortUrl))
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
	http.Redirect(w, r, "https://www."+u.FullUrl, http.StatusPermanentRedirect)

}

func (a *Application) CheckUrlStatus() {
	allLinks := a.repo.GetAllLinks()
	wg := sync.WaitGroup{}
	accessibleLinks := UrlSlice{}
	inaccessibleLinks := UrlSlice{}
	for i := 0; i < len(allLinks); i++ {
		wg.Add(1)
		go func(id int, link string) {
			defer wg.Done()
			resp, err := http.Get("http://" + link)
			if err != nil || resp.StatusCode != http.StatusOK {
				log.Println(err)
				if allLinks[id].Accessible == true {
					inaccessibleLinks.Lock()
					inaccessibleLinks.idSlice = append(inaccessibleLinks.idSlice, allLinks[id].Model.ID)
					inaccessibleLinks.Unlock()
				}
			} else {
				if allLinks[id].Accessible == false {
					accessibleLinks.Lock()
					accessibleLinks.idSlice = append(accessibleLinks.idSlice, allLinks[id].Model.ID)
					accessibleLinks.Unlock()
				}
			}
		}(i, allLinks[i].FullUrl)
	}
	wg.Wait()
	a.repo.UpdateAccess(inaccessibleLinks.idSlice, accessibleLinks.idSlice)
}

func (a *Application) GetShortUrlStats(w http.ResponseWriter, r *http.Request) {
	u := &models.Link{ShortUrl: r.URL.Query().Get("short")}
	err := a.repo.GetStats(u)
	if err != nil {
		w.Write([]byte("cant get stats for this short uri"))
		log.Println(err)
		return
	}
	stats := struct {
		FullUrl        string `json:"full_url"`
		NumOfRedirects uint32 `json:"number_of_redirects"`
		Accessible     bool   `json:"access_status"`
	}{
		FullUrl:        u.FullUrl,
		NumOfRedirects: u.NumsOfRedirects,
		Accessible:     u.Accessible,
	}
	answer, err := json.Marshal(stats)
	println(answer)
	if err != nil {
		w.Write([]byte("cant convert stats to json for this short uri"))
		log.Println(err)
		return
	}

	w.Write(answer)
}
