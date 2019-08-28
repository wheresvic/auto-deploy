package adserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
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

// ScriptResult ...
type ScriptResult struct {
	OutputSuccess string
	OutputError   string
}

// InitServer ...
func InitServer(initConfig *adconfiguration.AdConfiguration, adVersion adversion.AdVersion) *http.Server {

	port := strconv.Itoa(initConfig.Server.HTTPPort)

	r := mux.NewRouter()

	routerAPI := r.PathPrefix("/api").Subrouter()
	routerAPI.Use(loggingMiddleware)
	routerAPI.Use(corsMiddleware)

	// api
	routerAPI.HandleFunc("/version", wrapAPIHandler(apiHandler(getApiHandlerVersion(adVersion))))

	for _, project := range initConfig.Projects {

		// TODO: check for errors here

		route := "/webhooks/" + project.ProjectSlug

		if project.SCMServiceType != "" {
			route += "/" + project.SCMServiceType
		}

		log.Printf("%+v, %s", project, route)

		routerAPI.HandleFunc(route, wrapAPIHandler(apiHandler(getApiHandlerProjectSlug(project))))
	}

	fs := http.FileServer(http.Dir("public"))
	r.PathPrefix("/").Handler(fs)
	// r.Handle("/", fs)

	// log.Printf("%+v", routerAPI)

	server := &http.Server{Addr: ":" + port, Handler: r}

	// log.Println("Server listening on " + port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	// log.Fatal(http.ListenAndServe(":"+port, r))

	return server
}

func getApiHandlerVersion(adVersion adversion.AdVersion) func(w http.ResponseWriter, r *http.Request) *APIError {

	return func(w http.ResponseWriter, r *http.Request) *APIError {
		result, err := json.Marshal(adVersion)
		jsonAPIError := getAPIError(err)
		if jsonAPIError != nil {
			return jsonAPIError
		}
		fmt.Fprintf(w, string(result))
		return nil
	}

}

func getApiHandlerProjectSlug(project adconfiguration.AdProjectConfiguration) func(w http.ResponseWriter, r *http.Request) *APIError {

	return func(w http.ResponseWriter, r *http.Request) *APIError {
		var request interface{}
		// TODO: for a GET request this here is EOF
		err1 := json.NewDecoder(r.Body).Decode(&request)
		decodeJSONRequestBodyAPIError := getAPIError(err1)
		if decodeJSONRequestBodyAPIError != nil {
			return decodeJSONRequestBodyAPIError
		}

		log.Println("a")

		s, err2 := json.MarshalIndent(request, "", "\t")
		encodeJSONRequestBodyAPIError := getAPIError(err2)
		if encodeJSONRequestBodyAPIError != nil {
			return encodeJSONRequestBodyAPIError
		}

		// log.Printf("%+v", *request);
		log.Println(string(s))

		projectCommand := exec.Command(project.ProjectScript)

		projectCommandResult := ScriptResult{}

		projectCommandOutput, err := projectCommand.Output()
		if err != nil {
			projectCommandResult.OutputError = err.Error()
		} else {
			projectCommandResult.OutputSuccess = string(projectCommandOutput)
		}

		response, err3 := json.MarshalIndent(projectCommandResult, "", "\t")
		encodeJSONProcessCommandResultError := getAPIError(err3)
		if encodeJSONProcessCommandResultError != nil {
			return encodeJSONProcessCommandResultError
		}

		_, writeResponseError := w.Write(response)
		if writeResponseError != nil {
			log.Fatal(writeResponseError)
		}

		/*
			requestStrings := request.(map[string]interface{})

			if project.SCMServiceType == "github" {
				log.Println(requestStrings["ref"].(string))
			}
		*/

		// TODO: execute script and return results

		return nil

		/*
			result, err := json.Marshal(project)
			jsonAPIError := getAPIError(err)
			if jsonAPIError != nil {
				return jsonAPIError
			}
			fmt.Fprintf(w, string(result))
			return nil
		*/
	}
}

// Start ...
func Start(server *http.Server, port int) {
	log.Println("Server listening on " + strconv.Itoa(port))

	// returns ErrServerClosed on graceful close
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %s", err)
	}
}

func getAPIError(err error) *APIError {
	if err != nil {
		return &APIError{Error: err, ErrorMessage: err.Error(), Code: 500}
	}
	return nil
}
