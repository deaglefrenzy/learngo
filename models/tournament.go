package models

import (
	"time"
)

type Tournament struct {
	ID        string    `bson:"id,omitempty" json:"id,omitempty"`
	MatchID   []Match   `bson:"MatchID" json:"matchID"`
	Winner    string    `bson:"winner" json:"winner"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	DeletedAt time.Time `bson:"deletedAt" json:"deletedAt"`
}
