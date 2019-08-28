package adserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/wheresvic/auto-deploy/src/adconfiguration"
	"github.com/wheresvic/auto-deploy/src/adversion"
)

// ScriptResult ...
type ScriptResult struct {
	OutputSuccess string
	OutputError   string
}

func apiHandlerVersion(adVersion adversion.AdVersion) func(w http.ResponseWriter, r *http.Request) *APIError {

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

func apiHandlerProjectSlug(project adconfiguration.AdProjectConfiguration) func(w http.ResponseWriter, r *http.Request) *APIError {

	return func(w http.ResponseWriter, r *http.Request) *APIError {
		var request interface{}
		// TODO: for a GET request this here is EOF

		// first decode json payload
		err1 := json.NewDecoder(r.Body).Decode(&request)
		decodeJSONRequestBodyAPIError := getAPIError(err1)
		if decodeJSONRequestBodyAPIError != nil {
			return decodeJSONRequestBodyAPIError
		}

		// re-encode json payload for debugging
		s, err2 := json.MarshalIndent(request, "", "\t")
		encodeJSONRequestBodyAPIError := getAPIError(err2)
		if encodeJSONRequestBodyAPIError != nil {
			return encodeJSONRequestBodyAPIError
		}

		// log.Printf("%+v", *request);
		log.Println(string(s))

		go goExecuteProjectCommand(project)

		/*
			_, writeResponseError := w.Write(response)
			if writeResponseError != nil {
				log.Fatal(writeResponseError)
			}
		*/

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

func goExecuteProjectCommand(project adconfiguration.AdProjectConfiguration) {

	projectCommand := exec.Command(project.ProjectScript)
	projectCommand.Dir = project.ProjectRoot
	// projectCommand.Run()

	projectCommandResult := ScriptResult{}

	projectCommandOutput, err := projectCommand.Output()
	if err != nil {
		projectCommandResult.OutputError = err.Error()
	} else {
		projectCommandResult.OutputSuccess = string(projectCommandOutput)
	}

	projectCommandResultJSON, err3 := json.MarshalIndent(projectCommandResult, "", "\t")
	if err3 != nil {
		log.Fatal(err3)
	}

	/*
		encodeJSONProcessCommandResultError := getAPIError(err3)
		if encodeJSONProcessCommandResultError != nil {
			return encodeJSONProcessCommandResultError
		}
	*/

	log.Println(string(projectCommandResultJSON))
}

// Start ...
