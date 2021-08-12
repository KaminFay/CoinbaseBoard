package main

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

func initDB(dbUser string, dbPass string, dbHost string, dbPort string) *mongo.Client {
	mongoUri := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)
	log.Printf("Connecting to: mongodb://%s:****@%s:%s", dbUser, dbHost, dbPort)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to mongo db.")
	return client
}

func insertSingleDocIntoCollection(value Mongo_Document, coinCode string) {
	collection := CLIENT.Database("CoinbaseBoard").Collection(coinCode)
	result, err := collection.InsertOne(ctx, value)
	if err != nil {
		log.Println("There was an error inserting the document: ", err.Error())
	}
	log.Println("One value inserted:", reflect.TypeOf(value), result)
}

// func insertMultipleDocsIntoCollection(client *mongo.Client, value string) {

// }

func cleanupDB(client *mongo.Client) {
	client.Disconnect(ctx)
	log.Println("Disconnected the client")
}
