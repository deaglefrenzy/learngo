package handlers

import (
	"go_tutorial/models"
	"go_tutorial/repository"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CharHandler struct {
	repo *repository.MongoRepository
}

func NewCharHandler(repo *repository.MongoRepository) *CharHandler {
	return &CharHandler{repo: repo}
}

func (h *CharHandler) CreateChar(c *gin.Context) {

	var input struct {
		Name  string `json:"name"`
		Class string `json:"class"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.Name == "" || input.Class == "" {
		c.JSON(400, gin.H{"error": "Character name & class required."})
		return
	}
	id := h.GetAutoIncrement(c)

	newChar, err := models.NewCharacter(id, input.Name, input.Class)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err = h.repo.Create(newChar); err != nil {
		c.JSON(500, gin.H{"error": "Failed to insert character into MongoDB"})
		return
	}

	c.JSON(201, newChar)
}

func (h *CharHandler) GetAutoIncrement(c *gin.Context) int {

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
		c.JSON(500, gin.H{"error": "Failed to find auto increment"})
		return 0
	}

	update := bson.M{"$inc": bson.M{"value": 1}}

	err = sequenceRepo.UpdateOne(filter, update)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update auto increment"})
		return 0
	}

	return result.Value
}

func (h *CharHandler) IndexChars(c *gin.Context) {

	var characters []models.Character
	err := h.repo.FindMany(map[string]interface{}{}, &characters)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed fetch characters"})
		return
	}

	if characters == nil {
		characters = []models.Character{}
	}

	c.JSON(201, characters)
}

func (h *CharHandler) ShowChar(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.M{"id": id}
	var character models.Character

	err = h.repo.FindOne(filter, &character)
	if err != nil {
		c.JSON(404, gin.H{"error": "Character not found"})
		return
	}

	c.JSON(200, character)
}

func (h *CharHandler) UpdateName(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	var input struct {
		NewName string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if input.NewName == "" {
		c.JSON(400, gin.H{"error": "New character name required."})
		return
	}

	filter := bson.M{"basestatus.name": input.NewName}
	count, err := h.repo.CountDocuments(filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to check character uniqueness"})
		return
	}

	if count > 0 {
		c.JSON(409, gin.H{"error": "Character name already exists"})
		return
	}

	filter = bson.M{"id": id}
	update := bson.M{"$set": bson.M{"basestatus.name": input.NewName}}

	if err := h.repo.UpdateOne(filter, update); err != nil {
		c.JSON(500, gin.H{"error": "Failed to update character"})
		return
	}

	var character models.Character
	if err := h.repo.FindOne(filter, &character); err != nil {
		c.JSON(404, gin.H{"error": "Character not found"})
		return
	}

	c.JSON(200, character)
}

// type DeleteCharacterRequest struct {
// 	ID       string `uri:"id"`
// 	Password string `json:"password"`
// }

func (h *CharHandler) DestroyChar(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.M{"id": id}
	err = h.repo.DeleteOne(filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete character"})
		return
	}

	c.Status(204)
}

func (h *CharHandler) LevelUpChar(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	var character models.Character
	filter := bson.M{"id": id}
	err = h.repo.FindOne(filter, &character)
	if err != nil {
		c.JSON(404, gin.H{"error": "Character not found"})
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
		c.JSON(400, gin.H{"error": "Update level failed"})
		return
	}

	err = h.repo.FindOne(filter, &character)
	if err != nil {
		c.JSON(404, gin.H{"error": "Character not found"})
		return
	}

	c.JSON(200, character)
}
