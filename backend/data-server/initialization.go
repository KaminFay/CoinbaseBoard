package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func initService() *mongo.Client {
	log.Println("Initializing Services...")
	initializeEnvVars()
	getCoinbaseEnvVars()
	dbUser, dbPass, dbHost, dbPort := getDatabaseEnvVars()
	client := initDB(dbUser, dbPass, dbHost, dbPort)
	initializeCryptoList()
	log.Println("Services Complete")
	return client
}

func getDatabaseEnvVars() (string, string, string, string) {
	log.Println("Getting database env vars.")
	return os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT")
}

func getCoinbaseEnvVars() {
	log.Println("Getting Coinbase API env vars")
	API_KEY = os.Getenv("API_KEY")
	API_SECRET = os.Getenv("API_SECRET")
}

func initializeEnvVars() {
	log.Println("Loading .env files...")
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading environmental variables.")
	}
}

func initializeCryptoList() {
	log.Println("Loading coin list from Coinbase...")
	accountInfo := sendCoinbaseRequest("GET", "/v2/accounts?limit=100")
	json.Unmarshal(accountInfo, &CB_LISTED_CRYPTOS)
	log.Println("Initialized list of", len(CB_LISTED_CRYPTOS.Currencies), "coins.")
}
