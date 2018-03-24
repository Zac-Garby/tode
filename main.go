package main

import (
	"fmt"
	"net/http"

	"github.com/Zac-Garby/tode/api"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	api.Register(r)

	fmt.Println("listening on :7000")
	http.Handle("/", r)
	http.ListenAndServe(":7000", nil)
}
