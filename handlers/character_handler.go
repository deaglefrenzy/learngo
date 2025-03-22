package handlers

import (
	"encoding/json"
	"go_tutorial/models"
	"go_tutorial/utils"
	"net/http"
)

func CreateChar(w http.ResponseWriter, r *http.Request) {
	var char models.Character

	err := json.NewDecoder(r.Body).Decode(&char)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if char.BaseStatus.Name == "" || char.Class == "" {
		http.Error(w, "Character name & class required.", http.StatusBadRequest)
		return
	}

	newChar := models.NewCharacter(char.BaseStatus.Name, char.Class)

	utils.RespondJSON(w, newChar, http.StatusCreated)
}
