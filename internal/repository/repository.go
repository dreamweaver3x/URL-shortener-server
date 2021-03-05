package repository

import (
	"SOKR/internal/models"
	"SOKR/internal/shorturl"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"time"
)

type LinksRepository struct {
	db *gorm.DB
}

func NewLinksRepository(db *gorm.DB) *LinksRepository {
	return &LinksRepository{db: db}
}

func (l *LinksRepository) Create(u *models.Link) (*models.Link, error) {
	tx := l.db.Begin()
	result := tx.Create(u)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	u.ShortUrl = shorturl.Encode(int(u.ID))

	result = tx.Save(u)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	result = tx.Commit()
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	return u, nil
}

func (l *LinksRepository) CheckForElemLong(u *models.Link) bool {
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "http://")
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "https://")
	u.FullUrl = strings.TrimPrefix(u.FullUrl, "www.")
	result := l.db.Where("full_url = ?", u.FullUrl).First(&u)
	if result.Error != nil {
		return false
	}
	return true
}

func (l *LinksRepository) CheckForElemShort(u *models.Link) bool {
	result := l.db.Where("short_url = ?", u.ShortUrl).First(&u)
	if result.Error != nil {
		return false
	}
	if u.ShortUrl == "" {
		return false
	}
	return true
}

func (l *LinksRepository) GetLongUrl(u *models.Link) (*models.Link, error) {
	result := l.db.Where("short_url = ?", u.ShortUrl).First(&u)
	if result.Error != nil {
		return nil, result.Error
	}
	u.NumsOfRedirects++
	result = l.db.Model(u).Where("short_url = ?", u.ShortUrl).Update("nums_of_redirects", u.NumsOfRedirects)
	if result.Error != nil {
		return nil, result.Error
	}
	return u, nil
}

func (l *LinksRepository) CheckUrlsStatus() string {
	u := &models.Link{}

	for {
		l.db.Last(u)
		for i := 1; i <= int(u.ID)/10+1; i++ {
			go func(id int) {
				for k := id - 10; k < id; k++ {
					link := models.Link{}
					result := l.db.Where("id = ?", k).First(&link)
					if result.Error != nil {
						log.Println(result.Error)
						continue
					}
					resp, err := http.Get("http://www." + link.FullUrl)

					if err != nil || resp.StatusCode != 200 {
						log.Println(err)
						l.db.Table("links").Where("id = ?", k).Update("accessible", false)
						log.Println(link.FullUrl, " SSILKA NE RABOTAET")
						continue
					}
					err = resp.Body.Close()
					if err != nil {
						log.Println(err)
					}
				}
			}(i * 10)
		}
		time.Sleep(time.Minute * 10)
	}
}
