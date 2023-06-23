package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	var PORT = os.Getenv("PORT")
	var DATA_DIR = os.Getenv("DATA_DIR")
	var ENV = os.Getenv("ENV")

	if ENV == "development" || DATA_DIR == "" {
		DATA_DIR = "./data/"
	}

	if PORT == "" {
		PORT = "8080"
	}

	log.Printf("Starting server with env vars:\n\nENV: %s\nDATA_DIR: %s\nPORT: %s\n\n", ENV, DATA_DIR, PORT)

	fileServer := http.FileServer(http.Dir(DATA_DIR))

	http.Handle("/", WithLogging(fileServer))

	http.HandleFunc("/summaries/upload", getUploadFn("summaries/"))
	http.HandleFunc("/notes/upload", getUploadFn("notes/"))

	log.Printf("Listening on port %s ...", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
