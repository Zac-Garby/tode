package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// The tode API is the core functionality of tode. It can be used directly via a HTTP request, or by using the web interface.
//
// It supports a number of routes:
//
//  - /api/query/{op}/{query}
//  - /api/query/{op}/{query}/{limit | "first"}
//  - /api/random
//  - /api/random/{number}
//  - /api/user/{#id | name}
//  - /api/eq/{id}
//  - /api/all/users
//  - /api/all/equations
//
// Each route returns its result in the JSON format. You can probably guess what they all do.
// In the first two, {op} is one of ~, =, or !, which mean roughly, contains, and doesn't contain, respectively.
// If a request encounters an error, it returns some JSON looking something like {"error": "what happened?"}, possibly
// with more information.

// Register registers the API routes on
func Register(r mux.Router) {
	r.HandleFunc("/api/query/{op:(?:~|=|!)}/{query}", handleQuery)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func handleQueryLimit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func handleRandom(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func handleRandomLimit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func handleEquation(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func handleAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}

func handleAllEquations(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `{"error": "not implemented"}`)
}
