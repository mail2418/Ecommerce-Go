package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var UserCollection *mongo.Collection = UserData(Client, "Users")
var ProductCollection *mongo.Collection = ProductData(Client, "Products")

func DBSetup() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOps := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOps)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil{
		log.Println("failed to connect mongoDB")
		return nil
	}
	fmt.Println("Successfully connected to mongoDB")
	return client
}

var Client *mongo.Client = DBSetup()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	userCollection := client.Database("Ecommerce").Collection(collectionName)
	return userCollection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	productCollection := client.Database("Ecommerce").Collection(collectionName)
	return productCollection
}
