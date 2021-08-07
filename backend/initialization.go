package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var CB_LISTED_CRYPTOS CB_Listed_Crypto
var REQUESTED_CRYPTOS []string

func initializeCryptoList() {
	accountInfo := sendCoinbaseRequest("GET", "/v2/accounts?limit=100")
	json.Unmarshal(accountInfo, &CB_LISTED_CRYPTOS)
}

func getWalletValue(coinCode string) {
	walletContent, _ := getWalletContent(coinCode)
	currentCoinValue, _ := getCurrentCoinValue(coinCode)

	currentWalletValue := currentCoinValue * walletContent
	fmt.Println("Wallet Content: " + fmt.Sprintf("%f", walletContent) + " Coin Value: " + fmt.Sprintf("%f", currentCoinValue) + " Wallet Value: " + fmt.Sprintf("%f", currentWalletValue))
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
