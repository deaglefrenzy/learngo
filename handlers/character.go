package handlers

import (
	"context"
	"encoding/json"
	"go_tutorial/models"
	"go_tutorial/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(w http.ResponseWriter, r *http.Request, client *mongo.Client) {

	collection := client.Database("testing").Collection("characters")

	var input struct {
		Name  string `json:"name"`
		Class string `json:"class"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Class == "" {
		http.Error(w, "Character name & class required.", http.StatusBadRequest)
		return
	}

	id := GetAutoIncement(client)

	newChar, err := models.NewCharacter(id, input.Name, input.Class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = collection.InsertOne(context.Background(), newChar)
	if err != nil {
		http.Error(w, "Failed to insert character into MongoDB", http.StatusInternalServerError)
		log.Println("MongoDB insert error:", err)
		return
	}

	utils.RespondJSON(w, newChar, http.StatusCreated)
}

func GetAutoIncement(client *mongo.Client) int {
	collection := client.Database("testing").Collection("sequence")

	objectID, err := primitive.ObjectIDFromHex("67e4e853c2f4028ee4f83f4a")
	if err != nil {
		log.Fatalf("invalid ObjectId: %w", err)
		return 0
	}
	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"value": 1}}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var result struct {
		Value int `bson:"value"`
	}

	err = collection.FindOneAndUpdate(context.Background(), filter, update, options).Decode(&result)
	if err != nil {
		log.Fatalf("Error getting next value: %v", err)
		return 0
	}

	return result.Value
}

func Index(w http.ResponseWriter, r *http.Request, client *mongo.Client) {

	collection := client.Database("testing").Collection("characters")
	ctx := context.Background()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch characters: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var char []models.Character
	if err = cursor.All(ctx, &char); err != nil {
		http.Error(w, "Failed to decode characters: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if char == nil {
		char = []models.Character{}
	}

	utils.RespondJSON(w, char, http.StatusOK)
}

func Show(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	collection := client.Database("testing").Collection("characters")
	var character models.Character
	filter := bson.M{"id": id}

	err = collection.FindOne(context.Background(), filter).Decode(&character)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Character not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to query MongoDB", http.StatusInternalServerError)
		log.Println("MongoDB findOne error:", err)
		return
	}

	utils.RespondJSON(w, character, http.StatusOK)
}

func UpdateName(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	var input struct {
		NewName string `json:"name"`
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.NewName == "" {
		http.Error(w, "New character name required.", http.StatusBadRequest)
		return
	}

	collection := client.Database("testing").Collection("characters")

	count, _ := collection.CountDocuments(context.Background(), bson.M{"baseStatus.name": input.NewName})
	if count > 0 {
		http.Error(w, "Character name already exists", http.StatusConflict)
		return
	}

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"baseStatus.name": input.NewName}}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedCharacter models.Character
	err = collection.FindOneAndUpdate(context.Background(), filter, update, options).Decode(&updatedCharacter)
	if err != nil {
		http.Error(w, "Update Name Failed", http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, updatedCharacter, http.StatusOK)
}

func Destroy(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	collection := client.Database("testing").Collection("characters")

	filter := bson.M{"id": id}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "Delete Failed.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
