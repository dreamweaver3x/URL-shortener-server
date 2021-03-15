# URL shortener server
This is my http web server with echo framework. It can shrink the URL and contain information about the number redirects and accessibility of the link, using the short URL.
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
func (a *Application) GetShortURL(c echo.Context) error 
```
It returns a short URL in JSON to client using POST method.

_____
```golang
func (a *Application) RedirectWithShortUrl(c echo.Context) error
```
It Redirects you to address, based on short URL using GET method.
So http://localhost:8080/aaaaI will get you to www.iwantthislinktobeshorter.com
_____
```golang
func (a *Application) GetShortUrlStats(c echo.Context) error 
```
This Hanlder returns information about URL in JSON using GET method.

____
```golang
func (a *Application) CheckUrlStatus()
```
This method checks every URL in databse for accessibility
### config
Just for config.

# how to run 
- Clone repository
- `docker-compose up -d` in console, it'll create db and server.
- Try to add new links in database
