package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Store struct {
	DB  *redis.Client
	CTX context.Context
}

// Write writes input to db
func (store *Store) Write(id string, input *Input) error {
	p, err := json.Marshal(input)
	if err != nil {
		return errors.Wrap(err, "error marshalling data")
	}

	err = store.DB.Set(store.CTX, id, p, 0).Err()
	if err != nil {
		return errors.Wrap(err, "error writing to db")
	}

	return nil
}

// Print prints inputed messages and deleting it from db
func (store *Store) Print(id string, input *Input) {
	t := time.Now()
	currentTime := t.Unix()

	delay := input.Time - currentTime
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)

	}

	fmt.Println(input.Msg)

	err := store.DB.Del(store.CTX, id).Err()
	if err != nil {
		log.Println(errors.Wrap(err, "error deleting from db"))
	}
}

// Restore scan db for left data, used in case of server restart
func (store *Store) Restore() {
	keys := store.DB.Keys(store.CTX, "*").Val()
	for _, key := range keys {
		b, err := store.DB.Get(store.CTX, key).Bytes()
		if err != nil {
			log.Println(errors.Wrap(err, "error restoring data"))
			continue
		}

		input := &Input{}
		err = json.Unmarshal(b, input)
		if err != nil {
			log.Println(errors.Wrap(err, "error unmarshalling data"))
			continue
		}

		store.Print(key, input)
	}
}
