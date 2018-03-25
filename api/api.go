package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// The API encapsulates the core functionality of tode. It can be used directly via a HTTP request, or by using the
// web interface.
//
// It supports a number of routes:
//
//  - /api/query/{op}/{query}
//  - /api/query/{op}/{query}/{limit | "first"}
//  - /api/random
//  - /api/random/{number}
//  - /api/user/{name}
//  - /api/user/id/{id}
//  - /api/eq/{id}
//  - /api/all/users
//  - /api/all/equations
//
// Each route returns its result in the JSON format. You can probably guess what they all do.
// In the first two, {op} is one of ~, =, !, or r, which mean roughly, contains, doesn't contain, and
// matches regex, respectively.
// If a request encounters an error, it returns some JSON looking something like
// {"error": "what happened?"}, possibly with more information.
type API struct {
	db *redis.Client
}

// Register registers the API routes on
func (a *API) Register(r *mux.Router) error {
	a.db = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := a.db.Ping().Result()
	if err != nil {
		return err
	}

	r.HandleFunc("/api/query/{op:(?:~|=|!|r)}/{query}", a.handleQuery)
	r.HandleFunc("/api/query/{op:(?:~|=|!|r)}/{query}/{limit:(?:[0-9]+|first)}", a.handleQueryLimit)
	r.HandleFunc("/api/random", a.handleRandom)
	r.HandleFunc("/api/random/{limit:[0-9]+}", a.handleRandomLimit)
	r.HandleFunc("/api/user/{name:[a-zA-Z0-9_-]+}", a.handleUser)
	r.HandleFunc("/api/user/id/{id:[0-9]+}", a.handleUserID)
	r.HandleFunc("/api/eq/{id:[0-9]+}", a.handleEquation)
	r.HandleFunc("/api/all/users", a.handleAllUsers)
	r.HandleFunc("/api/all/equations", a.handleAllEquations)

	return nil
}

func (a *API) handleQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func (a *API) handleQueryLimit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func (a *API) handleRandom(w http.ResponseWriter, r *http.Request) {
	equations, err := a.FetchRandomEquations(1)
	if err != nil {
		writeError(w, err)
		return
	}

	if len(equations) < 1 {
		writeError(w, ErrEquationNotExist)
		return
	}

	out, err := json.Marshal(equations[0])
	if err != nil {
		writeError(w, err)
		return
	}

	w.Write(out)
}

func (a *API) handleRandomLimit(w http.ResponseWriter, r *http.Request) {
	rawLimit := mux.Vars(r)["limit"]

	limit, err := strconv.ParseInt(rawLimit, 10, 64)
	if err != nil {
		writeError(w, err)
		return
	}

	equations, err := a.FetchRandomEquations(limit)
	if err != nil {
		writeError(w, err)
		return
	}

	out, err := json.Marshal(equations)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Write(out)
}

func (a *API) handleUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]

	user, err := a.FetchUserByName(name)
	if err != nil {
		writeError(w, err)
		return
	}

	out, err := json.Marshal(user)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Write(out)
}

func (a *API) handleUserID(w http.ResponseWriter, r *http.Request) {
	rawID := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		writeError(w, err)
		return
	}

	user, err := a.FetchUser(id)
	if err != nil {
		writeError(w, err)
		return
	}

	out, err := json.Marshal(user)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Write(out)
}

func (a *API) handleEquation(w http.ResponseWriter, r *http.Request) {
	rawID := mux.Vars(r)["id"]

	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		writeError(w, err)
		return
	}

	equation, err := a.FetchEquation(id)
	if err != nil {
		writeError(w, err)
		return
	}

	out, err := json.Marshal(equation)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Write(out)
}

func (a *API) handleAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.FetchAllUsers()
	if err != nil {
		writeError(w, err)
		return
	}

	out, err := json.Marshal(users)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Write(out)
}

func (a *API) handleAllEquations(w http.ResponseWriter, r *http.Request) {
	equations, err := a.FetchAllEquations()
	if err != nil {
		writeError(w, err)
		return
	}

	out, err := json.Marshal(equations)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Write(out)
}
