package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jasonlvhit/gocron"
	"go.mongodb.org/mongo-driver/mongo"
)

/* global */
var accessToken string
var API_KEY string = ""
var API_SECRET string = ""
var CLIENT *mongo.Client
var CB_LISTED_CRYPTOS CB_Listed_Crypto
var REQUESTED_CRYPTOS []string
var SHUTDOWN = false

func main() {
	CLIENT = initService()
	go startCoinbaseMonitoring()
	defer cleanupDB(CLIENT)
	http.HandleFunc("/list_cryptos", listCryptos)
	http.HandleFunc("/get_crypto_from_list", getCryptoFromList)
	http.HandleFunc("/getBitcoin", getBitcoinInfo)
	fmt.Println(">>>>>>> OClient started at:", 8080)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func startCoinbaseMonitoring() {
	gocron.Every(1).Minute().Do(testing)
	<-gocron.Start()
}

func testing() {
	fmt.Println("Testing a print:", time.Now())
}

func listCryptos(w http.ResponseWriter, r *http.Request) {
	availValues, _ := json.Marshal(CB_LISTED_CRYPTOS)
	fmt.Fprintf(w, string(availValues), nil)
}

func getCryptoFromList(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(reqBody, &REQUESTED_CRYPTOS)

	for _, x := range REQUESTED_CRYPTOS {
		getWalletValue(x)
	}
}

// func getSingleCrypto(w http.ResponseWriter, r *http.Request) {
// }

func getBitcoinInfo(w http.ResponseWriter, r *http.Request) {
	accountInfo := sendCoinbaseRequest("GET", "/v2/accounts/BTC")
	fmt.Fprintf(w, string(accountInfo), nil)
}
