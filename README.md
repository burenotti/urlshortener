# UrlShortener

**URLShortener** is a simple pet-project written
in go that implements basic url-shortener api.
The project uses **PostgreSQL** as a main storage
and a **Redis** as a cache.

### How to run

First, provide a `.env` file with following content:

```shell
POSTGRES_USER=user
POSTGRES_PASSWORD=pass
POSTGRES_DB=urlshortener

DB_DSN=postgresql://user:pass@localhost:5432/urlshortener
REDIS_ADDR=localhost:6379
REDIS_PASS=pass
REDIS_DB=0
```

Then you can use `docker-compose` to start the server:

```shell
# will start server on localhost:8080
docker-compose -f docker-compose.yml up
```

Swagger docs are available on
[http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html).
