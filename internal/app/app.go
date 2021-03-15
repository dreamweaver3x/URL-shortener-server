package app

import (
	"SOKR/internal/models"
	"SOKR/internal/repository"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	_ "net/url"
	_ "strings"
	"sync"
	"time"
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

func (a *Application) Start(port string) {
	go func() {
		for {
			a.CheckUrlStatusNew()
			time.Sleep(time.Minute * 10)
		}
	}()

	e := echo.New()
	e.GET("/shortstats", a.GetShortUrlStats)
	e.POST("/urlshortener", a.GetShortURL)
	e.GET("/:uri", a.RedirectWithShortUrl)
	e.Logger.Fatal(e.Start(port))
}

func (a *Application) GetShortURL(c echo.Context) error {
	req := &struct {
		FullUrl  string `json:"full_url"`
		ShortUrl string `json:"short_url"`
	}{}
	u := &models.Link{}
	if err := c.Bind(req); err != nil {
		log.Println(err)
		return err
	}
	u.FullUrl = req.FullUrl
	err := a.repo.Create(u)
	if err != nil {
		return err
	}
	req.ShortUrl = "http://localhost:8080/" + u.ShortUrl

	return c.JSON(http.StatusAccepted, req)
}

func (a *Application) RedirectWithShortUrl(c echo.Context) error {
	u := &models.Link{}
	u.ShortUrl = c.Param("uri")
	log.Println("short = ", u.ShortUrl)
	err := a.repo.GetLongUrl(u)
	if err != nil {
		return err
	}
	println(u.FullUrl)

	return c.Redirect(http.StatusMovedPermanently, u.FullUrl)

}

func (a *Application) GetShortUrlStats(c echo.Context) error {
	u := &models.Link{}
	req := &struct {
		FullUrl        string `json:"full_url"`
		ShortUrl       string `json:"short_url"`
		NumOfRedirects uint32 `json:"number_of_redirects"`
		Accessible     bool   `json:"access_status"`
	}{}
	if err := c.Bind(req); err != nil {
		return err
	}
	u.ShortUrl = req.ShortUrl
	if err := a.repo.GetStats(u); err != nil {
		return err
	}
	req.FullUrl = u.FullUrl
	req.Accessible = u.Accessible
	req.NumOfRedirects = u.NumsOfRedirects

	return c.JSON(http.StatusAccepted, req)
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
			log.Println(err)
			break
		}
	L:
		for {
			select {
			case x = <-ch:
				wg.Add(1)
				go func(link models.Link) {
					println(i)
					defer wg.Done()
					resp, err := http.Get(link.FullUrl)
					if err != nil || resp.StatusCode != http.StatusOK {
						log.Println(err)
						if err == nil {
							log.Println(resp.StatusCode, "  ", link.FullUrl)
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
							if err := resp.Body.Close(); err != nil {
								log.Println(err)
							}
						}
					}

				}(x)

			default:

				break L
			}
		}
		i++
	}
	wg.Wait()
	err := a.repo.UpdateAccess(inaccessibleLinks.idSlice, accessibleLinks.idSlice)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("checked all urls")
	}
}
