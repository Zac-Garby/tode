package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// The API encapsulates the core functionality of tode. It can be used directly via a HTTP request, or by using the
// web interface.
//
// It supports a number of routes:
//
//     GET  /api/query/{op}/{query}
//     GET  /api/query/{op}/{query}/{limit | "all"}
//     GET  /api/random
//     GET  /api/random/{number}
//     GET  /api/user/{name}
//     GET  /api/user/id/{id}
//     GET  /api/equation/{id}
//     GET  /api/all/users
//     GET  /api/all/equations
//
//     PUT  /api/equation
//     PUT  /api/user
//
//  DELETE  /api/equation
//  DELETE  /api/user
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
	addr := os.Getenv("REDIS")
	if addr == "" {
		addr = "localhost:6379"
	}

	pw := os.Getenv("REDIS_PW")
	if pw == "" {
		pw = ""
	}

	db := os.Getenv("REDIS_DB")
	if db == "" {
		db = "0"
	}

	dbi, err := strconv.ParseInt(db, 10, 32)
	if err != nil {
		return fmt.Errorf("env var $REDIS_DB (%s) is not an integer", db)
	}

	a.db = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       int(dbi),
	})

	_, err = a.db.Ping().Result()
	if err != nil {
		return err
	}

	var (
		get = r.Methods("GET").Subrouter()
		put = r.Methods("PUT").Subrouter()
		del = r.Methods("DELETE").Subrouter()
	)

	get.HandleFunc("/api/query/{op:~|=|!|r}/{query}", a.handleQuery)
	get.HandleFunc("/api/query/{op:~|=|!|r}/{query}/{limit:[0-9]+|all}", a.handleQueryLimit)
	get.HandleFunc("/api/random", a.handleRandom)
	get.HandleFunc("/api/random/{limit:[0-9]+}", a.handleRandomLimit)
	get.HandleFunc("/api/user/{name:[a-zA-Z0-9_-]+}", a.handleUser)
	get.HandleFunc("/api/user/id/{id:[0-9]+}", a.handleUserID)
	get.HandleFunc("/api/equation/{id:[0-9]+}", a.handleEquation)
	get.HandleFunc("/api/all/users", a.handleAllUsers)
	get.HandleFunc("/api/all/equations", a.handleAllEquations)

	put.HandleFunc("/api/equation", a.handlePutEquation)
	put.HandleFunc("/api/user", a.handlePutUser)

	del.HandleFunc("/api/equation", a.handleDeleteEquation)
	del.HandleFunc("/api/user", a.handleDeleteUser)

	return nil
}

func (a *API) handleQuery(w http.ResponseWriter, r *http.Request) {
	var (
		op    = mux.Vars(r)["op"]
		query = mux.Vars(r)["query"]
		limit = int64(1)
		qt    QueryType
	)

	switch op {
	case "=":
		qt = QueryContainExact
	case "!":
		qt = QueryNotContain
	case "r":
		qt = QueryRegex
	default:
		qt = QueryContain
	}

	equations, err := a.Query(query, qt, limit)
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

func (a *API) handleQueryLimit(w http.ResponseWriter, r *http.Request) {
	var (
		op       = mux.Vars(r)["op"]
		query    = mux.Vars(r)["query"]
		rawLimit = mux.Vars(r)["limit"]
		limit    int64
		qt       QueryType
	)

	if rawLimit == "all" {
		limit = 0x7FFFFFFFFFFFFFFF // max int64
	} else {
		lim, err := strconv.ParseInt(rawLimit, 10, 64)
		if err != nil {
			writeError(w, err)
			return
		}
		limit = lim
	}

	switch op {
	case "=":
		qt = QueryContainExact
	case "!":
		qt = QueryNotContain
	case "r":
		qt = QueryRegex
	default:
		qt = QueryContain
	}

	equations, err := a.Query(query, qt, limit)
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

func (a *API) handlePutEquation(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func (a *API) handlePutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func (a *API) handleDeleteEquation(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func (a *API) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}
