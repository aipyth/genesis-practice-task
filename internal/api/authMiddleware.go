package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type TokenOnly struct {
	Token string `json:"token"`
}

func (s *Server) EnsureAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dataJson := TokenOnly{}
		if err := json.Unmarshal(data, &dataJson); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenData, err := s.storage.Get(dataJson.Token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if tokenData == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		email := tokenData.([]string)[0]

		ctx := context.WithValue(r.Context(), "user", email)
		r = r.WithContext(ctx)
		r.Body = ioutil.NopCloser(bytes.NewReader(data))

		h.ServeHTTP(w, r)
	}
}
