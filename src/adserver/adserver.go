package adserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wheresvic/auto-deploy/src/adconfiguration"
	"github.com/wheresvic/auto-deploy/src/adversion"
)

// APIError ...
type APIError struct {
	Error        error
	ErrorMessage string
	Code         int
}

// InitServer ...
func InitServer(initConfig *adconfiguration.AdConfiguration, adVersion adversion.AdVersion) {

	port := strconv.Itoa(initConfig.Server.HTTPPort)

	r := mux.NewRouter()

	routerAPI := r.PathPrefix("/api").Subrouter()
	routerAPI.Use(loggingMiddleware)
	routerAPI.Use(corsMiddleware)

	// api
	routerAPI.HandleFunc("/version", wrapAPIHandler(apiHandler(func(w http.ResponseWriter, r *http.Request) *APIError {
		result, err := json.Marshal(adVersion)
		jsonAPIError := getAPIError(err)
		if jsonAPIError != nil {
			return jsonAPIError
		}
		fmt.Fprintf(w, string(result))
		return nil
	})))

	for _, project := range initConfig.Projects {
		route := "/webhooks/" + project.ProjectSlug

		log.Printf("%+v, %s", project, route)

		routerAPI.HandleFunc(route, wrapAPIHandler(apiHandler(func(w http.ResponseWriter, r *http.Request) *APIError {
			result, err := json.Marshal(project)
			jsonAPIError := getAPIError(err)
			if jsonAPIError != nil {
				return jsonAPIError
			}
			fmt.Fprintf(w, string(result))
			return nil
		})))
	}

	log.Printf("%+v", routerAPI)

	log.Println("Server listening on " + port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getAPIError(err error) *APIError {
	if err != nil {
		return &APIError{Error: err, ErrorMessage: err.Error(), Code: 500}
	}
	return nil
}
