package app

import (
	"SOKR/internal/models"
	"SOKR/internal/repository"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/url"
	"strings"
	_ "strings"
	"sync"
)

type UrlSlice struct {
	sync.Mutex
	idSlice []uint
}

type ShortLinkRequest struct {
	FullUrl string `json:"full_url"`
	ShortUrl string `json:"short_url"`
}

type Application struct {
	repo *repository.LinksRepository
}

func NewApplication(repo *repository.LinksRepository) *Application {
	return &Application{repo: repo}
}

func (a *Application) GetShortURL(w http.ResponseWriter, r *http.Request) {
	req := &ShortLinkRequest{}
	u := &models.Link{}
	body, err := ioutil.ReadAll(r.Body)
	log.Println(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u.FullUrl = req.FullUrl
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "http://")
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "https://")
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "www.")
	err = a.repo.Create(u)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req.ShortUrl = "http://localhost:8080/" + u.ShortUrl
	answer, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Write(answer)
	w.WriteHeader(http.StatusOK)
}

func (a *Application) RedirectWithShortUrl(w http.ResponseWriter, r *http.Request) {
	u := &models.Link{}
	linkWithSlash := r.URL.String()
	u.ShortUrl = linkWithSlash[1:]
	u, err := a.repo.GetLongUrl(u)
	if err != nil {
		w.Write([]byte("this short url doesn't work"))
		log.Println(err)
		return
	}
	println(u.FullUrl)
	http.Redirect(w, r, "http://"+u.FullUrl, http.StatusPermanentRedirect)

}



func (a *Application) CheckUrlStatusNew() {
	wg := sync.WaitGroup{}
	ch := make(chan models.Link, 500)
	accessibleLinks := UrlSlice{}
	inaccessibleLinks := UrlSlice{}
	var i uint
	var x models.Link
	for {
		err := a.repo.GetFiveHundredLinks(i, &ch)
		if err != nil {
			log.Println("NU VOT BLIN", err)
			break
		}

	L:
		for {
			select {
			case x = <-ch:
				wg.Add(1)
				go func(link models.Link) {
					//println(x.FullUrl)
					defer wg.Done()
					resp, err := http.Get("http://" + link.FullUrl)
					if err != nil || resp.StatusCode != http.StatusOK {
						//log.Println(err)
						if err == nil {
							//log.Println(resp.StatusCode, "  ", link.FullUrl)
						}
						if link.Accessible == true {
							inaccessibleLinks.Lock()
							inaccessibleLinks.idSlice = append(inaccessibleLinks.idSlice, link.Model.ID)
							inaccessibleLinks.Unlock()
						}
					} else {
						if link.Accessible == false {
							accessibleLinks.Lock()
							accessibleLinks.idSlice = append(accessibleLinks.idSlice, link.Model.ID)
							accessibleLinks.Unlock()
							resp.Body.Close()
						}
					}
				}(x)

			default:

				break L
			}

		}
		println("YA TUT VSEM KU")
		i++
	}
	a.repo.UpdateAccess(inaccessibleLinks.idSlice, accessibleLinks.idSlice)
	wg.Wait()
	log.Println("zaconchil smotret' ssilki")
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
