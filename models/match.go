package models

import "time"

type Match struct {
	ID        string      `bson:"_id,omitempty" json:"id,omitempty"`
	TeamA     []Character `bson:"teamA" json:"teamA"`
	TeamB     []Character `bson:"teamB" json:"teamB"`
	CreatedAt time.Time   `bson:"createdAt" json:"createdAt"`
	Winner    string      `bson:"winner,omitempty" json:"winner,omitempty"`
}
