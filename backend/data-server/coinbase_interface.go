package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var CACHED_TIMESTAMP = ""

/**
 * getTimeStamp() string is used to get a timestamp from the coinbase api for use with authentication
 * as the server may not be in sync with the users local machine. This value is then returned.
 *
 * Note: This call does not need authentication.
 */
func getTimeStamp() {
	resp, err := http.Get("https://api.coinbase.com/v2/time")
	if err != nil {
		log.Error("Issue making the call to coinbase /v2/time")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Issue reading body response: ", err.Error())
	}

	var timestampData CB_TimeStampData
	json.Unmarshal(body, &timestampData)
	CACHED_TIMESTAMP = strconv.FormatInt(timestampData.Data.Epoch, 10)
}

/**
 * generateSignature(message, secret string) generates a hmac string that represents the required authentication message
 * for a coinbase API call and returns it.
 */
func generateSignature(message, secret string) (string, error) {
	sha := sha256.New
	hash := hmac.New(sha, []byte(secret))
	hash.Write([]byte(message))
	signature := fmt.Sprintf("%x", hash.Sum(nil))
	return signature, nil
}

func createSignature(method, uri string) (string, error) {
	getTimeStamp()
	message := fmt.Sprintf("%s%s%s", CACHED_TIMESTAMP, method, uri)
	return message, nil
}

func sendCoinbaseRequest(method, uri string) []byte {
	baseURI := "https://api.coinbase.com"
	authMessage, err := createSignature(method, uri)
	if err != nil {
		log.Error("Issue create auth message: ", err.Error())
	}

	signature, err := generateSignature(authMessage, API_SECRET)
	if err != nil {
		log.Error("Issues generating signature: ", err.Error())
	}

	client := http.Client{}
	req, err := http.NewRequest(method, baseURI+uri, nil)
	if err != nil {
		log.Error("Error creating new http request: ", err.Error())
	}

	req.Header = http.Header{
		"CB-ACCESS-SIGN":      []string{signature},
		"CB-ACCESS-TIMESTAMP": []string{CACHED_TIMESTAMP},
		"CB-ACCESS-KEY":       []string{API_KEY},
		"Content-Type":        []string{"application/json"},
		"CB-VERSION":          []string{"2021-07-25"},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("There was an error sending request: ", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response bytes: ", err)
	}

	return body
}

func getWalletValue(coinCode string) {
	walletContent, _ := getWalletContent(coinCode)
	currentCoinValue, _ := getCurrentCoinValue(coinCode)

	currentWalletValue := currentCoinValue * walletContent
	fmt.Println("Wallet Content: " + fmt.Sprintf("%f", walletContent) + " Coin Value: " + fmt.Sprintf("%f", currentCoinValue) + " Wallet Value: " + fmt.Sprintf("%f", currentWalletValue))
	valueForMongo := Mongo_Document{walletContent,
		currentCoinValue,
		currentWalletValue,
		float64(time.Now().Unix())}
	insertSingleDocIntoCollection(valueForMongo, coinCode)
}

func getCurrentCoinValue(coinCode string) (float64, error) {
	spotValueData := sendCoinbaseRequest("GET", "/v2/prices/"+coinCode+"-USD/spot")
	spotValue := CB_Spot_Price{}
	json.Unmarshal(spotValueData, &spotValue)
	return strconv.ParseFloat(spotValue.Data.Amount, 64)
}

func getWalletContent(coinCode string) (float64, error) {
	walletValueData := sendCoinbaseRequest("GET", "/v2/accounts/"+coinCode)
	singleCrypto := CB_Single_Crypto{}
	json.Unmarshal(walletValueData, &singleCrypto)
	return strconv.ParseFloat(singleCrypto.Data.Balance.Amount, 64)
}
