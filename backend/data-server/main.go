package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
)

/* global */
var accessToken string
var API_KEY string = ""
var API_SECRET string = ""

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func initializeEnvVars() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading environmental variables.")
	}

	API_KEY = os.Getenv("API_KEY")
	API_SECRET = os.Getenv("API_SECRET")
}

func main() {
	initLogger()
	initializeEnvVars()
	initializeCryptoList()
	fmt.Println("Testing")
	http.HandleFunc("/list_cryptos", listCryptos)
	http.HandleFunc("/get_crypto_from_list", getCryptoFromList)
	http.HandleFunc("/getBitcoin", getBitcoinInfo)
	fmt.Println(">>>>>>> OClient started at:", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func listCryptos(w http.ResponseWriter, r *http.Request) {
	availValues, _ := json.Marshal(CB_LISTED_CRYPTOS)
	fmt.Fprintf(w, string(availValues), nil)
}

func getCryptoFromList(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	var testing []string
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(reqBody, &testing)

	for _, x := range testing {
		getWalletValue(x)
	}
}

func getBitcoinInfo(w http.ResponseWriter, r *http.Request) {
	accountInfo := sendCoinbaseRequest("GET", "/v2/accounts/BTC")
	fmt.Fprintf(w, string(accountInfo), nil)
}
