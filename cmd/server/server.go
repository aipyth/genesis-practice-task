package server

import (
	"github.com/aipyth/genesis-practice-task/internal/api"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/user/create", api.CreateUser).Methods("POST")
	router.HandleFunc("/user/login", api.LoginUser).Methods("POST")
	router.HandleFunc("/btcRate", api.EnsureAuth(api.GetBTCRateInUAH)).Methods("GET")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServe()
}
