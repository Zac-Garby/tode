package api

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrDatabase = errors.New("database error")

	ErrUserNotExist         = errors.New("user doesn't exist")
	ErrUserInvalidTimestamp = errors.New("user has non-int timestamp")
	ErrUserInvalidID        = errors.New("user has invalid ID")

	ErrEquationNotExist         = errors.New("equation doesn't exist")
	ErrEquationInvalidTimestamp = errors.New("equation has non-int timestamp")
	ErrEquationInvalidScore     = errors.New("equation has non-int score")
	ErrEquationInvalidAuthor    = errors.New("equation has non-int author")
)

func writeError(w http.ResponseWriter, e error) {
	http.Error(w, fmt.Sprintf(`{"error": "%s"}`, e.Error()), http.StatusInternalServerError)
}
