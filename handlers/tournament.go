package handlers

import (
	"fmt"
	"go_tutorial/models"
	"go_tutorial/repository"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type TournamentHandler struct {
	repo *repository.MongoRepository
}

func NewTournamentHandler(repo *repository.MongoRepository) *TournamentHandler {
	return &TournamentHandler{repo: repo}
}

func (h *TournamentHandler) StartTournament(c *gin.Context) {

	matchRepo := h.repo.WithCollection("matches")

	var inputs struct {
		MatchID []string `json:"matchID"`
	}

	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	seen := make(map[string]bool)
	for _, id := range inputs.MatchID {
		if seen[id] {
			c.JSON(400, gin.H{"error": "duplicate matchID found: " + id})
			return
		}
		seen[id] = true
	}

	var matches []models.Match
	for _, id := range inputs.MatchID {
		var match models.Match
		err := matchRepo.FindOne(bson.M{"id": id}, &match)
		if err != nil {
			c.JSON(404, gin.H{"error": "Match not found: " + id})
			return
		}
		matches = append(matches, match)
	}

	var win []string
	for _, m := range matches {
		fmt.Println(time.Now())
		var myCh = make(chan string, 1)
		go func() {
			myCh <- models.Battle(m)
		}()
		win = append(win, <-myCh)
		fmt.Println(time.Now())
	}
	fmt.Println(win)

	// tournament := models.Tournament{
	// 	ID:        uuid.NewString(),
	// 	MatchID:   matches,
	// 	CreatedAt: time.Now(),
	// }

	// if err := h.repo.Create(&tournament); err != nil {
	// 	c.JSON(500, gin.H{"error": "Failed to save tournament: " + err.Error()})
	// 	return
	// }

	// c.JSON(200, tournament)
	c.Status(204)
}
