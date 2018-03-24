package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Zac-Garby/tode/api"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	a := new(api.API)

	if err := a.Register(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("listening on http://localhost:7000")
	http.Handle("/", r)
	http.ListenAndServe(":7000", nil)
}
