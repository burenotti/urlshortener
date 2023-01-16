package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/burenotti/urlshortener/internal/handler"
	"github.com/burenotti/urlshortener/internal/service"
	"github.com/burenotti/urlshortener/internal/storage"
	sessions "github.com/gin-contrib/sessions/redis"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	var host string
	var port int
	flag.StringVar(&host, "host", "localhost", "start server on given host")
	flag.IntVar(&port, "port", 8080, "start server on given port")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", host, port)

	err := InitConfig(".env")
	pool, err := pgxpool.New(context.TODO(), viper.GetString("DB_DSN"))
	err = pool.Ping(context.TODO())

	if err != nil {
		logrus.Fatalf("can't connect to postgres: %s", err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_ADDR"),
		Password: viper.GetString("REDIS_PASS"),
		DB:       viper.GetInt("REDIS_DB"),
	})
	_, err = rdb.Ping(context.TODO()).Result()

	if err != nil {
		logrus.Fatalf("can't connect to redis: %s", err.Error())
	}

	mainShortener := storage.NewPostgresStorage(pool)
	cacheShortener := storage.NewRedisStorage(rdb, 1*time.Hour)
	composedShortener := storage.NewComposedShortener(mainShortener, cacheShortener)
	statStore := storage.NewPgStatStore(pool)
	store := storage.NewStorage(composedShortener, statStore)

	serv := service.NewService(store)
	viper.SetDefault("BASE_PATH", fmt.Sprintf("http://%s/l", addr))

	sessionStore, err := sessions.NewStore(
		10,
		"tcp",
		viper.GetString("REDIS_ADDR"),
		viper.GetString("REDIS_PASS"),
		[]byte(viper.GetString("SESSION_AUTH_KEY")),
	)

	handle := handler.NewHandler(serv, sessionStore, viper.GetString("BASE_PATH"))

	logrus.Infof("starting server on %s", addr)

	err = http.ListenAndServe(addr, handle)
	if err != nil {
		logrus.Error("serving error: %s", err.Error())
	}

}

func InitConfig(file string) error {
	viper.SetConfigFile(file)
	return viper.ReadInConfig()
}
