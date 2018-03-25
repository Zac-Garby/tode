package api

import (
	"fmt"
	"strconv"
)

// An Equation is a representation of an equation from the database.
type Equation struct {
	ID          int64    `json:"id"`
	Source      string   `json:"source"`
	Description string   `json:"description"`
	Author      int64    `json:"author"`
	Categories  []string `json:"categories"`
	Score       int64    `json:"score"`
	Confirmed   bool     `json:"confirmed"`
	Timestamp   int64    `json:"added"`
}

// FetchEquation fetches an equation from the database, by id.
func (a *API) FetchEquation(id int64) (*Equation, error) {
	key := fmt.Sprintf("equation:%d", id)

	val, err := a.db.HGetAll(key).Result()
	if err != nil {
		return nil, ErrDatabase
	}

	if len(val) == 0 {
		return nil, ErrEquationNotExist
	}

	author, err := strconv.ParseInt(val["author"], 10, 64)
	if err != nil {
		return nil, ErrEquationInvalidAuthor
	}

	score, err := strconv.ParseInt(val["score"], 10, 64)
	if err != nil {
		return nil, ErrEquationInvalidScore
	}

	timestamp, err := strconv.ParseInt(val["timestamp"], 10, 64)
	if err != nil {
		return nil, ErrUserInvalidTimestamp
	}

	categories, err := a.getCategories(id)
	if err != nil {
		return nil, err
	}

	return &Equation{
		ID:          id,
		Source:      val["source"],
		Description: val["description"],
		Author:      author,
		Categories:  categories,
		Score:       score,
		Confirmed:   val["confirmed"] == "yes",
		Timestamp:   timestamp,
	}, nil
}

// FetchAllEquations creates a slice containing all equations in the database,
// calling FetchEquation for each one.
func (a *API) FetchAllEquations() ([]*Equation, error) {
	ids, err := a.db.SMembers("equations").Result()
	if err != nil {
		return nil, ErrDatabase
	}

	equations := make([]*Equation, len(ids))

	for i, idStr := range ids {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, ErrUserInvalidID
		}

		equation, err := a.FetchEquation(id)
		if err != nil {
			return nil, ErrEquationNotExist
		}

		equations[i] = equation
	}

	return equations, nil
}

// FetchRandomEquations fetches n random equations from the database.
func (a *API) FetchRandomEquations(n int64) ([]*Equation, error) {
	ids, err := a.db.SRandMemberN("equations", n).Result()
	if err != nil {
		return nil, ErrDatabase
	}

	equations := make([]*Equation, len(ids))

	for i, idStr := range ids {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, ErrEquationInvalidID
		}

		eq, err := a.FetchEquation(id)
		if err != nil {
			return nil, err
		}

		equations[i] = eq
	}

	return equations, nil
}

func (a *API) getCategories(id int64) ([]string, error) {
	key := fmt.Sprintf("equation:%d:categories", id)

	val, err := a.db.SMembers(key).Result()
	if err != nil {
		return nil, ErrDatabase
	}

	return val, nil
}
