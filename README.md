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
It returns a short URL in JSON to client using POST method. Using cURL, the request will look like this:
```
curl -X POST http://localhost:8080/urlshortener \
> -H 'Content-Type: application/json' \
> -d '{"full_url":"https://www.youtube.com/watch?v=7oEZaljP7uY&t=3s"}'
```
And the answer: 
```
{"full_url":"https://www.youtube.com/watch?v=7oEZaljP7uY\u0026t=3s","short_url":"http://localhost:8080/aaaGd"}
```

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
This Hanlder returns information about URL in JSON using GET method. Using cURL, it will look like this:
```
curl -X GET http://localhost:8080/shortstats \
-H 'Content-Type: application/json' \
-d '{"short_url":"aaaGf"}'
```
And the answer:
```
{"full_url":"https://www.youtube.com/watch?v=7oEZaljP7uY\u0026t=3s","short_url":"aaaGf","number_of_redirects":3,"access_status":true}
```

____
```golang
func (a *Application) CheckUrlStatus()
```
This method checks every URL in databse for accessibility
```golang
func (a *Application) CheckUrlStatusNew()
```
In this method we check accessibility of our URLs. Ut takes 500 links from DB, put them into channel, goroutines will listen to channel, use GET on link, if code status wont be `200` and access status of the link is `true`, it will add ID of the link in a slice to update accesibility in DB later. If there are more links then 500, it will repeat until it takes all the links. 500 links is acceptable to accupy our memory.
### config
Just for config.

# how to run 
- Clone repository
- `docker-compose up -d` in console, it'll create db and server.
- Try to add new links in database
