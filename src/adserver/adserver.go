package adserver

import (
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
func InitServer(initConfig *adconfiguration.AdConfiguration, adVersion adversion.AdVersion) *http.Server {

	port := strconv.Itoa(initConfig.Server.HTTPPort)

	r := getRouter(initConfig, adVersion)

	server := &http.Server{Addr: ":" + port, Handler: r}

	// log.Println("Server listening on " + port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	// log.Fatal(http.ListenAndServe(":"+port, r))

	return server
}

// getRouter ...
func getRouter(initConfig *adconfiguration.AdConfiguration, adVersion adversion.AdVersion) *mux.Router {

	r := mux.NewRouter()

	routerAPI := r.PathPrefix("/api").Subrouter()
	routerAPI.Use(loggingMiddleware)
	routerAPI.Use(corsMiddleware)

	// api
	routerAPI.HandleFunc("/version", wrapAPIHandler(apiHandler(apiHandlerVersion(adVersion))))

	for _, project := range initConfig.Projects {

		// TODO: check for errors here

		route := "/webhooks/" + project.ProjectSlug

		if project.SCMServiceType != "" {
			route += "/" + project.SCMServiceType
		}

		log.Printf("%+v, %s", project, route)

		routerAPI.HandleFunc(route, wrapAPIHandler(apiHandler(apiHandlerProjectSlug(project))))
	}

	fs := http.FileServer(http.Dir("public"))
	r.PathPrefix("/").Handler(fs)
	// r.Handle("/", fs)

	// log.Printf("%+v", routerAPI)

	return r
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
