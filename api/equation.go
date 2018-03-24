package api

import (
	"fmt"

	"github.com/go-redis/redis"
)

// An Equation is a representation of an equation from the database.
type Equation struct {
	ID          int64    `json:"id"`
	Source      string   `json:"source"`
	Description string   `json:"description"`
	Author      int64    `json:"author"`
	Categories  []string `json:"categories"`
	Score       int      `json:"score"`
	Confirmed   bool     `json:"confirmed"`
	Timestamp   int64    `json:"added"`
}

// FetchEquation fetches an equation from the database, by id.
func FetchEquation(db *redis.Client, id int64) (*Equation, error) {
	key := fmt.Sprintf("equation:%d", id)

	val, err := db.HGetAll(key).Result()
	if err != nil {
		return nil, ErrDatabase
	}

	if len(val) == 0 {
		return nil, ErrEquationNotExist
	}
}
