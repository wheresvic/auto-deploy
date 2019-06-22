package adversion

import (
	"testing"
)

// TestUnitGetCurrentVersion
func TestUnitGetCurrentVersion(t *testing.T) {
	version := GetCurrentVersion()

	expectedVersionGo := "go1.12.1"

	if version.VersionGo != expectedVersionGo {
		t.Errorf("Failed, expected %s, got %s", expectedVersionGo, version.VersionGo)
	}
}
