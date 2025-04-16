package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_tutorial/db"
	"go_tutorial/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindManyCharacterRequest struct {
	Class string
}

func main() {
	client, err := db.ConnectMongoDB()
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return
	}
	defer db.CloseMongoDB(client)
	db := client.Database("testing")

	query := bson.D{
		{Key: "class", Value: "paladin"},
		{Key: "basestatus.health", Value: bson.M{"$gt": 70}},
	}

	coll := db.Collection("characters")

	cur, err := coll.Find(context.TODO(), query, options.Find().SetLimit(20))
	if err != nil {
		panic(err)
	}

	var characters []models.Character
	if err := cur.All(context.TODO(), &characters); err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(characters)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))
}
