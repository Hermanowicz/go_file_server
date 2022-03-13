package main

// TODO: replace ListenAndServe with ListenAndServeTLS
// TODO: add makefile with cert setup and build command
// generating ssl -> https://devcenter.heroku.com/articles/ssl-certificate-self

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	// default location which will be served
	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal("Error wille riding cwd, terminating app.")
	}

	// default port
	port := flag.String("port", "5000", "Port to use for file server.")
	dir := flag.String("dir", cwd, "directory to serve. Default current working dir.")

	// info for user
	fmt.Printf("Starting file server on PORT: %s with DIR: %s", *port, *dir)

	// meaty part of this program
	handler := http.FileServer(http.Dir(string(*dir)))
	log.Fatal(http.ListenAndServe(":"+*port, handler))
}
