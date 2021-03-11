# URL shortener server
This is my http web server. It can shrink the URL and contain information about the number redirects and accessibility of the link, using the short URL.
## packages
### repository
Here's where all actions with database happen.
### models
There's struct for GORM + automigrate in case you want to add new columns to your database.
### shorturl
It makes a new short URL based on id in the database. For the 1st link it return "aaaaa", for the 2nd "aaaab" and so on. Pretty sure it's not the best way, but it's working =)
### app
All 3 handlers are in this package. 
```golang
func (a *Application) GetShortURL(w http.ResponseWriter, r *http.Request)
```
It returns a short URL to client
If you want a new short URl, type
[http:/localhost:8080/urlshortener?link=www.iwantthislinktobeshorter.com](http:/localhost:8080/urlshortener?link=www.iwantthislinktobeshorter.com)
You'll get something like http://localhost:8080/aaaaI
_____
```golang
func (a *Application) RedirectWithShortUrl(w http.ResponseWriter, r *http.Request)
```
It Redirects you to address, based on short URL
So http://localhost:8080/aaaaI (from previous func) will get you to www.iwantthislinktobeshorter.com
_____
```golang
func (a *Application) GetShortUrlStats(w http.ResponseWriter, r *http.Request)
```
This Hanlder returns information about URL
example: for http://localhost:8080/getshortstats?short=aaaag you'll get {"full_url":"iwantthislinktobeshorter.com","number_of_redirects":0,"access_status":true}
____
```golang
func (a *Application) CheckUrlStatus()
```
This method checks every URL in databse for accessibility


# how to run 
- Clone repository
- Create database. You can use this command: `docker-compose exec pgdb psql -U db_user -c 'CREATE DATABASE URLcutter'`
- Run the main file (check config in main.go, so your connection to database will be successful), you can use commands in Makefile
- Try to add new links in database
