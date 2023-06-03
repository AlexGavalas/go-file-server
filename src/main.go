package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func getUploadFn(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Uploading file ... ")

		// 10 << 20 specifies a max upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)

		// FormFile returns the first file for the given key
		// it also returns the FileHeader so we can get
		// the Filename and the size of the file
		file, handler, err := r.FormFile("md_file")

		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}

		defer file.Close()

		fmt.Printf("File: %v\n", handler.Filename)
		fmt.Printf("File size: %.2fKB\n", float32(handler.Size)/float32(1000))

		var DATA_DIR = os.Getenv("DATA_DIR")
		var ENV = os.Getenv("ENV")

		if ENV == "dev" || DATA_DIR == "" {
			DATA_DIR = "./data/"
		}

		DATA_DIR = DATA_DIR + prefix

		mkdirErr := os.MkdirAll(DATA_DIR, os.ModePerm)

		if mkdirErr != nil {
			log.Println(err)
		}

		// Create a temporary file within the directory that follows
		// a particular naming pattern
		tempFile, err := os.CreateTemp(DATA_DIR, "upload-*.md")

		if err != nil {
			fmt.Println(err)
		}

		defer tempFile.Close()

		log.Println("Created temp file: " + tempFile.Name())

		// read all of the contents of our uploaded file into a byte array
		fileBytes, err := io.ReadAll(file)

		if err != nil {
			fmt.Println(err)
		}

		// write this byte array to our temporary file
		tempFile.Write(fileBytes)

		os.Rename(tempFile.Name(), DATA_DIR+handler.Filename)

		log.Println("Success renamed file!")

		fmt.Fprintf(w, "[OK]\nSuccessfully Uploaded File\n")
	}
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {

		uri := r.RequestURI
		method := r.Method
		h.ServeHTTP(rw, r) // serve the original request

		// log request details
		log.Println("Request details", map[string]interface{}{
			"uri":    uri,
			"method": method,
		})
	}

	return http.HandlerFunc(logFn)
}

func main() {
	var PORT = os.Getenv("PORT")
	var DATA_DIR = os.Getenv("DATA_DIR")
	var ENV = os.Getenv("ENV")

	log.Println(("Starting server with env vars ..."))
	log.Println("ENV: " + ENV)
	log.Println("DATA_DIR: " + DATA_DIR)
	log.Println("PORT: " + PORT)

	if ENV == "dev" || DATA_DIR == "" {
		DATA_DIR = "./data/"
	}

	if PORT == "" {
		PORT = "8080"
	}

	fileServer := http.FileServer(http.Dir(DATA_DIR))

	http.Handle("/", WithLogging(fileServer))

	http.HandleFunc("/summaries/upload", getUploadFn("summaries/"))

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello there!\n")
	})

	log.Printf("Listening on port %s ...", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, nil))
}
