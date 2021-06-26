package api

import (
	"github.com/aipyth/genesis-practice-task/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Server struct {
	storage 		storage.Storage
	router 			*mux.Router
	Addr 			string
	ReadTimeout		time.Duration
	WriteTimeout 	time.Duration
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)

		duration := time.Since(start)
		log.Println(method, uri, duration)
	})
}

func (s *Server) InitRoutes() {
	s.router.Use(logger)

	s.router.HandleFunc("/user/create", s.CreateUser).Methods("POST")
	s.router.HandleFunc("/user/login", s.LoginUser).Methods("POST")
	s.router.HandleFunc("/btcRate", s.EnsureAuth(s.GetBTCRateInUAH)).Methods(
		"GET")
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr: s.Addr,
		Handler: s.router,
		WriteTimeout: s.WriteTimeout,
		ReadTimeout: s.ReadTimeout,
	}
	return server.ListenAndServe()
}



func Run() error {
	serv := Server{
		storage: 		storage.NewCSVStorage("/tmp/genesis-api-task"),
		router:  		mux.NewRouter(),
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	if err := serv.storage.Connect(); err != nil {
		return err
	}

	serv.InitRoutes()

	return serv.Start()
}
