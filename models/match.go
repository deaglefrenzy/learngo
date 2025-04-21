package models

import (
	"math/rand"
	"time"
)

type Match struct {
	ID        string      `bson:"id,omitempty" json:"id,omitempty"`
	TeamA     []Character `bson:"teamA" json:"teamA"`
	TeamB     []Character `bson:"teamB" json:"teamB"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt"`
	Winner    string      `bson:"winner" json:"winner"`
}

func CharAttackDefense(c Character) (Character, int, int) {
	att := c.BaseStatus.Attack
	def := 0
	newShield := c.Shield
	newMana := c.Mana
	switch c.Class {
	case "paladin":
		if newShield > 0 {
			att += newShield
			def += newShield
			newShield -= 1
		}
	case "archer":
		if c.Critical > 0 && rand.Intn(100) <= c.Critical {
			att *= 2
		}
		if c.Critical > 0 && rand.Intn(100) <= c.Critical {
			def += 2
		}
	case "mage":
		if newMana >= 7 {
			att += 7
			newMana -= 7
		}
		if newMana >= 3 {
			def += 3
			newMana -= 3
		}
	}
	char := Character{
		ID: c.ID,
		BaseStatus: BaseStatus{
			Name:   c.BaseStatus.Name,
			Health: c.BaseStatus.Health,
			Attack: c.BaseStatus.Attack,
		},
		Level:    c.Level,
		Class:    c.Class,
		Shield:   newShield,
		Critical: c.Critical,
		Mana:     newMana,
	}
	return char, att, def
}
