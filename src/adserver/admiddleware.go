package adserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rs/cors"
)

type apiHandler func(http.ResponseWriter, *http.Request) *APIError

var c = cors.AllowAll()

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {

		log.Printf("http %s %s | %+v", r.Method, r.URL.Path, e.Error)

		response, err := json.Marshal(e)
		// will this ever happen?
		if err != nil {
			http.Error(w, e.Error.Error(), e.Code)
		}
		// happy case
		http.Error(w, string(response), e.Code)
	}
}

func wrapAPIHandler(wrappedAPIHandler apiHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// c.Handler(wrappedAPIHandler) // the cors handler will call the ServeHTTP method
		w.Header().Set("Content-Type", "application/json")
		wrappedAPIHandler.ServeHTTP(w, req)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if (*req).Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, req)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("http %s %s", req.Method, req.URL.Path)
		next.ServeHTTP(w, req)
	})
}
