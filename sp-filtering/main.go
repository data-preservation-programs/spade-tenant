package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	godotenv.Load()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	fmt.Println("starting...")
	fmt.Println(os.Getenv("MONGODB_URI"))
	fmt.Println("end...")
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	groupStage := bson.D{{"$sort", bson.D{"data_time", "-1"}},
		{"$group", bson.D{
			{"_id", "$category"},
			{"average_price", bson.D{{"$avg", "$price"}}},
			{"type_total", bson.D{{"$sum", 1}}},
		}}}

	// Send a ping to confirm a successful connection
	kentiks := client.Database("reputation").Collection("kentiks")
	var kentik Kentik
	err = kentiks.FindOne(context.TODO(), bson.D{{}}).Decode(&kentik)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(kentik)
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}
