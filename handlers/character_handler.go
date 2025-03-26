package handlers

import (
	"encoding/json"
	"go_tutorial/models"
	"go_tutorial/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateChar(w http.ResponseWriter, r *http.Request) {
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

	newChar, err := models.NewCharacter(input.Name, input.Class)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := models.LoadArrayChar("characters.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonData = append(jsonData, newChar)
	models.CharToJSON("characters.json", jsonData)

	utils.RespondJSON(w, newChar, http.StatusCreated)
}

func IndexChars(w http.ResponseWriter, r *http.Request) {
	filename := "characters.json"
	filepath := filepath.Join("database", filename)
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		http.Error(w, "error loading JSON file.", http.StatusBadRequest)
		return
	}

	var char []models.Character
	err = json.Unmarshal(fileData, &char)
	if err != nil {
		http.Error(w, "error unmarshaling JSON data", http.StatusBadRequest)
		return
	}

	utils.RespondJSON(w, char, http.StatusCreated)
}

func ShowChar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}
	filename := "characters.json"
	filepath := filepath.Join("database", filename)
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		http.Error(w, "Error loading JSON file.", http.StatusInternalServerError)
		return
	}

	var characters []models.Character
	err = json.Unmarshal(fileData, &characters)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON data", http.StatusInternalServerError)
		return
	}

	for _, char := range characters {
		if char.ID == id {
			utils.RespondJSON(w, char, http.StatusOK)
			return
		}
	}
	http.Error(w, "Character not found", http.StatusNotFound)
}

func ChangeCharName(w http.ResponseWriter, r *http.Request) {
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

	characters, err := models.LoadArrayChar("characters.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, char := range characters {
		if char.BaseStatus.Name == input.NewName && char.ID != id {
			http.Error(w, "Character name already exists.", http.StatusBadRequest)
			return
		}
	}

	updatedCharacters := make([]models.Character, 0, len(characters))
	for _, char := range characters {
		if char.ID == id {
			char.BaseStatus.Name = input.NewName
		}
		updatedCharacters = append(updatedCharacters, char)
	}
	models.CharToJSON("characters.json", updatedCharacters)
	utils.RespondJSON(w, updatedCharacters[id-1], http.StatusOK)
}

func DeleteChar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid character ID", http.StatusBadRequest)
		return
	}

	characters, err := models.LoadArrayChar("characters.json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedCharacters := make([]models.Character, 0, len(characters))
	for _, char := range characters {
		if char.ID != id {
			updatedCharacters = append(updatedCharacters, char)
		}
	}
	models.CharToJSON("characters.json", updatedCharacters)
	w.WriteHeader(http.StatusOK)
}
