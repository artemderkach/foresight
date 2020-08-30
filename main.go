package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Env struct {
	APP_PORT string
	DB_PORT  string
}

func main() {
	rest := &Rest{}
	env := parseEnv()

	// app dependencies
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:" + env.DB_PORT,
		Password: "",
		DB:       0,
	})
	rest.Store = &Store{
		DB:  rdb,
		CTX: context.Background(),
	}

	// print messages that left in db
	rest.Store.Restore()

	fmt.Println("runnig server on", env.APP_PORT)
	err := http.ListenAndServe(":"+env.APP_PORT, rest.Router())
	if err != nil {
		log.Fatal(errors.Wrap(err, "error starting server"))
	}
}

func parseEnv() *Env {
	env := &Env{}

	port, exists := os.LookupEnv("APP_PORT")
	if !exists {
		port = "8080"
	}
	env.APP_PORT = port

	dbPort, exists := os.LookupEnv("DB_PORT")
	if !exists {
		port = "6379"
	}
	env.DB_PORT = dbPort

	return env
}
