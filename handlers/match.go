package handlers

import (
	"errors"
	"fmt"
	"go_tutorial/models"
	"go_tutorial/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MatchHandler struct {
	repo *repository.MongoRepository
}

func NewMatchHandler(repo *repository.MongoRepository) *MatchHandler {
	return &MatchHandler{repo: repo}
}

func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var inputs struct {
		A []int `json:"A"`
		B []int `json:"B"`
	}

	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(inputs.A) != len(inputs.B) {
		c.JSON(400, gin.H{"error": "Teams need to be the same size"})
		return
	}

	seen := make(map[int]bool)
	for _, id := range append(inputs.A, inputs.B...) {
		if seen[id] {
			c.JSON(400, gin.H{"error": fmt.Sprintf("Character ID %d is duplicated", id)})
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
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		matchup.TeamA = append(matchup.TeamA, char)
	}

	for _, id := range inputs.B {
		char, err := h.getValidatedChar(id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
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

	if err := h.repo.Create(&match); err != nil {
		c.JSON(500, gin.H{"error": "Failed to save match: " + err.Error()})
		return
	}

	c.JSON(200, match)
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

func (h *MatchHandler) GetMatch(c *gin.Context) {
	matchID := c.Param("id")

	if matchID == "" {
		c.JSON(400, gin.H{"error": "Match ID is required"})
		return
	}

	match, err := h.getMatchByID(matchID)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, match)
}

func (h *MatchHandler) DestroyMatch(c *gin.Context) {
	matchID := c.Param("id")

	if matchID == "" {
		c.JSON(400, gin.H{"error": "Match ID is required"})
		return
	}

	filter := bson.M{"id": matchID}
	err := h.repo.DeleteOne(filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete match: " + err.Error()})
		return
	}

	c.Status(204)
}

func (h *MatchHandler) StartBattle(c *gin.Context) {
	matchID := c.Param("id")

	if matchID == "" {
		c.JSON(400, gin.H{"error": "Match ID is required"})
		return
	}

	match, err := h.getMatchByID(matchID)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	winner := models.Battle(match)

	fmt.Println(winner)
	filter := bson.M{"id": matchID}
	update := bson.M{"$set": bson.M{"winner": winner}}

	err = h.repo.UpdateOne(filter, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to write match winner: " + err.Error()})
		return
	}

	c.JSON(200, winner)
}
