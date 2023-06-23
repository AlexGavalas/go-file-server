package main

import (
	"log"
	"net/http"
)

func WithLogging(h http.Handler) http.Handler {
	logFn := func(rw http.ResponseWriter, r *http.Request) {

		uri := r.RequestURI
		method := r.Method

		// serve the original request
		h.ServeHTTP(rw, r)

		log.Println("Request details", map[string]interface{}{
			"uri":    uri,
			"method": method,
		})
	}

	return http.HandlerFunc(logFn)
}
