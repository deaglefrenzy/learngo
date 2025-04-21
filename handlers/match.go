package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_tutorial/models"
	"go_tutorial/repository"
	"go_tutorial/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
		char, err := h.getValidatedChar(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Use the error returned by the helper
			return
		}
		matchup.TeamA = append(matchup.TeamA, char)
	}

	for _, id := range inputs.B {
		char, err := h.getValidatedChar(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // Use the error returned by the helper
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

	utils.RespondJSON(w, match, http.StatusOK)
}

func (h *MatchHandler) getValidatedChar(id int) (models.Character, error) {
	var char models.Character
	filter := bson.M{"id": id}
	err := h.repo.WithCollection("characters").FindOne(filter, &char)
	if err != nil {
		return char, fmt.Errorf("character with ID %d not found", id)
	}

	filter = bson.M{
		"$or": []bson.M{
			{"teamA.id": id},
			{"teamB.id": id},
		},
		"winner": "",
	}
	var existingMatch models.Match
	err = h.repo.FindOne(filter, &existingMatch)
	if err == nil {
		return char, fmt.Errorf("character %d is already in an ongoing match: %s", id, existingMatch.ID)
	}
	return char, nil
}

func (h *MatchHandler) getMatchByID(matchID string) (models.Match, error) {
	var match models.Match
	err := h.repo.FindOne(bson.M{"id": matchID}, &match)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.Match{}, fmt.Errorf("match with ID %s not found", matchID)
		}
		return models.Match{}, fmt.Errorf("error fetching match with ID %s: %w", matchID, err) //wrap the error
	}
	return match, nil
}

func (h *MatchHandler) GetMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchID := vars["id"]

	if matchID == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}

	match, err := h.getMatchByID(matchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	utils.RespondJSON(w, match, http.StatusOK)
}

func (h *MatchHandler) DestroyMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchID := vars["id"]

	if matchID == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}

	filter := bson.M{"id": matchID}
	err := h.repo.DeleteOne(filter)
	if err != nil {
		http.Error(w, "Failed to delete match: "+err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MatchHandler) StartBattle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matchID := vars["id"]

	if matchID == "" {
		http.Error(w, "Match ID is required", http.StatusBadRequest)
		return
	}

	match, err := h.getMatchByID(matchID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	teamA := match.TeamA
	teamB := match.TeamB

	totalHealthA := 0
	totalHealthB := 0

	for _, c := range teamA {
		totalHealthA += c.BaseStatus.Health
	}
	for _, c := range teamB {
		totalHealthB += c.BaseStatus.Health
	}

	dead := false
	result := ""
	winner := ""

	for !dead {
		attackA := 0
		attackB := 0
		defenseA := 0
		defenseB := 0
		att := 0
		def := 0
		var newCharStat models.Character

		for i := 0; i < len(teamA); i++ {
			newCharStat, att, def = models.CharAttackDefense(teamA[i])
			teamA[i] = newCharStat
			attackA += att
			defenseA += def
		}
		for i := 0; i < len(teamA); i++ {
			newCharStat, att, def = models.CharAttackDefense(teamB[i])
			teamB[i] = newCharStat
			attackB += att
			defenseB += def
		}

		damageA := attackB - defenseA
		damageB := attackA - defenseB
		fmt.Printf("Team A att: %d, def: %d, dmg: %d\n", attackA, defenseA, damageA)
		fmt.Printf("Team B att: %d, def: %d, dmg: %d\n", attackB, defenseB, damageB)
		totalHealthA -= damageA
		totalHealthB -= damageB
		fmt.Printf("Team A's HP: %d\n", totalHealthA)
		fmt.Printf("Team B's HP: %d\n", totalHealthB)
		if totalHealthA <= 0 {
			dead = true
			if totalHealthB <= 0 {
				if totalHealthA < totalHealthB {
					winner = "B"
				} else {
					winner = "A"
				}
			} else {
				result = "TEAM B WIN"
				winner = "B"
			}
		} else if totalHealthB <= 0 {
			dead = true
			result = "TEAM A WIN"
			winner = "A"
		}
	}
	fmt.Println(result)
	filter := bson.M{"id": matchID}
	update := bson.M{"$set": bson.M{"winner": winner}}

	err = h.repo.UpdateOne(filter, update)
	if err != nil {
		http.Error(w, "Failed to write match winner", http.StatusInternalServerError)
		return
	}

	utils.RespondJSON(w, result, http.StatusOK)
}
