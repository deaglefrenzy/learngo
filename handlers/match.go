package handlers

import (
	"encoding/json"
	"fmt"
	"go_tutorial/models"
	"go_tutorial/repository"
	"go_tutorial/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type MatchHandler struct {
	repo *repository.MongoRepository
}

func NewMatchHandler(repo *repository.MongoRepository) *MatchHandler {
	return &MatchHandler{repo: repo}
}

func (h *MatchHandler) CreateMatch(w http.ResponseWriter, r *http.Request) {

	var inputs struct {
		A []int `json:"A"`
		B []int `json:"B"`
	}

	err := json.NewDecoder(r.Body).Decode(&inputs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(inputs.A) != len(inputs.B) {
		http.Error(w, "Teams needs to be the same size", http.StatusBadRequest)
		return
	}

	seen := make(map[int]bool)
	for _, id := range append(inputs.A, inputs.B...) {
		if seen[id] {
			http.Error(w, fmt.Sprintf("Character ID %d is duplicated", id), http.StatusBadRequest)
			return
		}
		seen[id] = true
	}

	var matchup struct {
		TeamA []models.Character `json:"teamA"`
		TeamB []models.Character `json:"teamB"`
	}

	for _, id := range inputs.A {
		var char models.Character
		filter := bson.M{"id": id}
		err = h.repo.WithCollection("characters").FindOne(filter, &char)
		if err != nil {
			http.Error(w, fmt.Sprintf("Character with ID %d not found", id), http.StatusNotFound)
			return
		}
		matchup.TeamA = append(matchup.TeamA, char)
	}

	for _, id := range inputs.B {
		var char models.Character
		filter := bson.M{"id": id}
		err = h.repo.WithCollection("characters").FindOne(filter, &char)
		if err != nil {
			http.Error(w, fmt.Sprintf("Character with ID %d not found", id), http.StatusNotFound)
			return
		}
		matchup.TeamB = append(matchup.TeamB, char)
	}

	match := models.Match{
		ID:        uuid.NewString(),
		TeamA:     matchup.TeamA,
		TeamB:     matchup.TeamB,
		CreatedAt: time.Now(),
	}

	err = h.repo.Create(&match)
	if err != nil {
		http.Error(w, "Failed to save match: "+err.Error(), http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, matchup, http.StatusOK)
}
