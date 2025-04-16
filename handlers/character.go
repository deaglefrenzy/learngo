package handlers

import (
	"encoding/json"
	"fmt"
	"go_tutorial/models"
	"go_tutorial/repository"
	"go_tutorial/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharHandler struct {
	repo *repository.MongoRepository
}

func NewCharHandler(repo *repository.MongoRepository) *CharHandler {
	return &CharHandler{repo: repo}
}

func (h *CharHandler) CreateChar(w http.ResponseWriter, r *http.Request) {

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
	id := h.GetAutoIncrement(w, r)

	newChar, err := models.NewCharacter(id, input.Name, input.Class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.repo.Create(newChar); err != nil {
		http.Error(w, "Failed to insert character into MongoDB", http.StatusInternalServerError)
		log.Println("MongoDB insert error:", err)
		return
	}

	utils.RespondJSON(w, newChar, http.StatusCreated)
}

func (h *CharHandler) GetAutoIncrement(w http.ResponseWriter, r *http.Request) int {

	sequenceRepo := h.repo.WithCollection("sequence")

	objectID, err := primitive.ObjectIDFromHex("67e4e853c2f4028ee4f83f4a")
	if err != nil {
		log.Fatalf("invalid ObjectId: %s", err)
		return 0
	}
	filter := bson.M{"_id": objectID}

	var result struct {
		Value int `bson:"value"`
	}

	err = sequenceRepo.FindOne(filter, &result)
	if err != nil {
		log.Println("Auto increment find error:", err)
		http.Error(w, "Failed to find auto increment", http.StatusInternalServerError)
		return 0
	}

	update := bson.M{"$inc": bson.M{"value": 1}}

	err = sequenceRepo.UpdateOne(filter, update)
	if err != nil {
		http.Error(w, "Failed to update auto increment", http.StatusInternalServerError)
		return 0
	}

	return result.Value
}

func (h *CharHandler) IndexChars(w http.ResponseWriter, r *http.Request) {

	var characters []models.Character
	err := h.repo.FindMany(map[string]interface{}{}, &characters)
	if err != nil {
		http.Error(w, "Failed to fetch characters: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if characters == nil {
		characters = []models.Character{}
	}

	utils.RespondJSON(w, characters, http.StatusOK)
}

func (h *CharHandler) ShowChar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	filter := bson.M{"id": id}
	var character models.Character

	err = h.repo.FindOne(filter, &character)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(character)
}

func (h *CharHandler) UpdateName(w http.ResponseWriter, r *http.Request) {
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

	filter := bson.M{"basestatus.name": input.NewName}
	var character models.Character

	count, err := h.repo.CountDocuments(filter)
	if err != nil {
		http.Error(w, "Failed to check character uniqueness", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "Character name already exists", http.StatusConflict)
		return
	}

	fmt.Println(count)

	filter = bson.M{"id": id}
	update := bson.M{"$set": bson.M{"basestatus.name": input.NewName}}

	err = h.repo.UpdateOne(filter, update)
	if err != nil {
		http.Error(w, "Failed to update character", http.StatusInternalServerError)
		return
	}

	err = h.repo.FindOne(filter, &character)
	if err != nil {
		http.Error(w, "Character not found", http.StatusConflict)
		return
	}

	utils.RespondJSON(w, character, http.StatusOK)
}

// type DeleteCharacterRequest struct {
// 	ID       string `uri:"id"`
// 	Password string `json:"password"`
// }

func (h *CharHandler) DestroyChar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	filter := bson.M{"id": id}
	err = h.repo.DeleteOne(filter)
	if err != nil {
		http.Error(w, "Failed to delete character", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CharHandler) LevelUpChar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	var character models.Character
	filter := bson.M{"id": id}
	err = h.repo.FindOne(filter, &character)
	if err != nil {
		http.Error(w, "Character not found", http.StatusNotFound)
		return
	}

	newHealth := character.BaseStatus.Health + 5
	newAttack := character.BaseStatus.Attack + 1
	newLevel := character.Level + 1

	update := bson.M{"$set": bson.M{
		"basestatus.health": newHealth,
		"basestatus.attack": newAttack,
		"level":             newLevel,
	}}

	err = h.repo.UpdateOne(filter, update)
	if err != nil {
		http.Error(w, "Level Up Failed", http.StatusInternalServerError)
		return
	}

	err = h.repo.FindOne(filter, &character)
	if err != nil {
		http.Error(w, "Character not found", http.StatusConflict)
		return
	}

	utils.RespondJSON(w, character, http.StatusOK)
}
