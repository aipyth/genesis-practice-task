package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type btcRateInfo struct {
	Rate float64 `json:"rate"`
}

var apiKey = os.Getenv("API_KEY")
var btcRateApiUrl = "https://rest.coinapi." +
	"io/v1/exchangerate/BTC/UAH?apikey=" + apiKey


func GetBTCRateInUAH(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(btcRateApiUrl)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var responseJson btcRateInfo
	if err := json.Unmarshal(respBody, &responseJson); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	rate := fmt.Sprint(responseJson.Rate)
	_, err = w.Write([]byte(rate))
	if err != nil {
		log.Println(err)
	}
}
