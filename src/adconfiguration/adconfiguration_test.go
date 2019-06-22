package adconfiguration

import (
	"testing"
)

func TestUnitShouldThrowErrorOnGetConfigurationWhenNotLoaded(t *testing.T) {
	_, err := GetConfiguration()
	if err == nil {
		t.Errorf("Should have thrown error as configuration not loaded")
	}
}

func TestUnitShouldLoadAndSetConfiguration(t *testing.T) {
	initConfig, err := LoadAndSetConfiguration("../../config.json")
	if err != nil {
		t.Error(err)
	}

	if initConfig.Server.HTTPPort == 0 {
		t.Errorf("Should have loaded the http port %d", initConfig.Server.HTTPPort)
	}

	if len(initConfig.Projects) == 0 {
		t.Error("Should have loaded projects")
	}

	project := initConfig.Projects[0]

	if project.ProjectSlug == "" {
		t.Errorf("Should have loaded the project slug %s", project.ProjectSlug)
	}

	if project.ProjectRoot == "" {
		t.Errorf("Should have loaded the project root folder location %s", project.ProjectRoot)
	}

	if project.ProjectScript == "" {
		t.Errorf("Should have loaded the project deployment script location %s", project.ProjectScript)
	}
}
