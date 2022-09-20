package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = "8080"

type Config struct{}

func main() {
	app := Config{}

	log.Println("starting broker service on port:", PORT)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
