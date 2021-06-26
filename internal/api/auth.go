package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
)

const tokenLength = 10

type CreateUserParams struct {
	Email 		string `json:"email"`
	Password 	string `json:"password"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken() string {
	b := make([]byte, tokenLength)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", b)
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var requestJson CreateUserParams
	if err := json.Unmarshal(data, &requestJson); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//log.Println(requestJson)
	if requestJson.Email == "" || requestJson.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("email or password is not provided"))
		return
	}

	if r, _ := s.storage.Get(requestJson.Email); r != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	hashedPassword, err := hashPassword(requestJson.Password)
	if err != nil {
		log.Println(err)
	}
	token := generateToken()
	err = s.storage.Set(requestJson.Email, []string{hashedPassword, token})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := s.storage.Set(token, requestJson.Email); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	if err := s.storage.Save(); err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusCreated)
}

type LoginUserParams struct {
	Email		string	`json:"email"`
	Password	string 	`json:"password"`
}

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var requestJson LoginUserParams
	if err := json.Unmarshal(data, &requestJson); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	storedData, err := s.storage.Get(requestJson.Email)
	if err != nil {
		log.Println(err)
	}

	// theres no such user with this email
	if storedData == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	userData, ok := storedData.([]string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// password hash is first in list of data
	passMatch := checkPasswordHash(requestJson.Password, userData[0])

	if passMatch {
		// token is second in list of data
		response, err := json.Marshal(TokenOnly{Token: userData[1]})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(response)
		if err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}

}
