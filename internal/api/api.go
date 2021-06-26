package api

import (
	"github.com/aipyth/genesis-practice-task/internal/storage"
	"log"
)

var Storage = storage.NewCSVStorage("/tmp/genesis-api-task")

func init() {
	if err := Storage.Connect(); err != nil {
		log.Fatal(err)
	}
}
