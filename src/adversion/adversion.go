package adversion

import "runtime"

// please maintain line numbers here for proper versioning
const version = "1.0.0"
const lastUpdated = 1561233534

// AdVersion ...
type AdVersion struct {
	Version     string
	VersionGo   string
	LastUpdated int64
}

// GetCurrentVersion ...
func GetCurrentVersion() AdVersion {
	return AdVersion{version, runtime.Version(), lastUpdated}
}
