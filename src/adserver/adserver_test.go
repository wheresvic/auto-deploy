package adserver

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/wheresvic/auto-deploy/src/adconfiguration"
	"github.com/wheresvic/auto-deploy/src/adversion"
)

func setUp(t *testing.T) *http.Server {
	adVersion := adversion.GetCurrentVersion()
	initConfig, err := adconfiguration.LoadAndSetConfiguration("../../config.json")
	if err != nil {
		t.Fatal(err)
	}

	initConfig.Server.HTTPPort = 9112

	server := InitServer(initConfig, adVersion)
	go func() {
		Start(server, initConfig.Server.HTTPPort)
	}()

	return server
}

func tearDown(t *testing.T, server *http.Server) {
	if err := server.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}

// TestIntegrationGetVersion
func TestIntegrationGetVersion(t *testing.T) {
	server := setUp(t)
	response, err := http.Get("http://localhost:9112/api/version")
	if err != nil {
		t.Fatal(err)
	}
	var responseBody interface{}
	err1 := json.NewDecoder(response.Body).Decode(&responseBody)
	if err1 != nil {
		t.Fatal(err1)
	}

	/*
		s, err2 := json.MarshalIndent(responseBody, "", "\t")
		if err2 != nil {
			t.Fatal(err2)
		}
		log.Println(string(s))
	*/

	responseBodyStrings := responseBody.(map[string]interface{})

	expectedVersion := "1.0.1"
	actualVersion := responseBodyStrings["Version"].(string)

	if expectedVersion != actualVersion {
		t.Errorf("Expected %s, Actual %s", expectedVersion, actualVersion)
	}

	expectedVersionGo := "go1.12.1"
	actualVersionGo := responseBodyStrings["VersionGo"].(string)

	if expectedVersionGo != actualVersionGo {
		t.Errorf("Expected %s, Actual %s", expectedVersionGo, actualVersionGo)
	}

	tearDown(t, server)
}
