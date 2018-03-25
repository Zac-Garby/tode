package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Zac-Garby/tode/api"
	"github.com/gorilla/mux"
)

var port = os.Getenv("PORT")

func main() {
	if port == "" {
		port = "7000"
	}

	var (
		r = mux.NewRouter()
		a = new(api.API)
	)

	if err := a.Register(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("listening on http://localhost:" + port)
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}
