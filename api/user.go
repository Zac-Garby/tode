package api

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

// A User is an in-memory representation of a user from the database.
type User struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"-"`
	Hash      string  `json:"-"`
	Salt      string  `json:"-"`
	Equations []int64 `json:"equations"`
	Timestamp int64   `json:"joined"`
}

func FetchUser(db *redis.Client, id int64) (*User, error) {
	key := fmt.Sprintf("user:%d", id)

	val, err := db.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	timestamp, err := strconv.ParseInt(val["timestamp"], 10, 64)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Name:      val["name"],
		Email:     val["email"],
		Hash:      val["hash"],
		Salt:      val["salt"],
		Timestamp: timestamp,
	}, nil
}

func FetchUserByName(db *redis.Client, name string) (*User, error) {
	val, err := db.HGet("usernames", name).Result()
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, err
	}

	return FetchUser(db, id)
}
