package models

import (
	"math/rand"
	"time"
)

type Match struct {
	ID        string      `bson:"id,omitempty" json:"id,omitempty"`
	TeamA     []Character `bson:"teamA" json:"teamA"`
	TeamB     []Character `bson:"teamB" json:"teamB"`
	Winner    string      `bson:"winner" json:"winner"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt"`
	DeletedAt time.Time   `bson:"deletedAt" json:"deletedAt"`
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

func Battle(match Match) string {
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
	winner := ""

	for !dead {
		attackA := 0
		attackB := 0
		defenseA := 0
		defenseB := 0
		att := 0
		def := 0
		var newCharStat Character

		for i := 0; i < len(teamA); i++ {
			newCharStat, att, def = CharAttackDefense(teamA[i])
			teamA[i] = newCharStat
			attackA += att
			defenseA += def
		}
		for i := 0; i < len(teamA); i++ {
			newCharStat, att, def = CharAttackDefense(teamB[i])
			teamB[i] = newCharStat
			attackB += att
			defenseB += def
		}

		damageA := attackB - defenseA
		damageB := attackA - defenseB

		totalHealthA -= damageA
		totalHealthB -= damageB

		if totalHealthA <= 0 {
			dead = true
			if totalHealthB <= 0 {
				if totalHealthA < totalHealthB {
					winner = "B"
				} else {
					winner = "A"
				}
			} else {
				winner = "B"
			}
		} else if totalHealthB <= 0 {
			dead = true
			winner = "A"
		}
	}
	return winner
}
