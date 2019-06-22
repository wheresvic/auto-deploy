package adconfiguration

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

var config *AdConfiguration

// AdConfiguration ...
type AdConfiguration struct {
	Server struct {
		HTTPPort int
	}
	Projects []AdProjectConfiguration
}

// AdProjectConfiguration
type AdProjectConfiguration struct {
	ProjectSlug   string
	ProjectRoot   string
	ProjectScript string
}

// LoadAndSetConfiguration ...
func LoadAndSetConfiguration(path string) (*AdConfiguration, error) {
	var data AdConfiguration

	file, readFileErr := ioutil.ReadFile(path)
	if readFileErr != nil {
		return nil, errors.Wrap(readFileErr, "could not read configuration file")
	}

	jsonUnmarshallErr := json.Unmarshal([]byte(file), &data)
	if jsonUnmarshallErr != nil {
		return nil, errors.Wrap(readFileErr, "could not parse json into config format")
	}
	config = &data
	return config, nil
}

// GetConfiguration ...
func GetConfiguration() (*AdConfiguration, error) {
	if config == nil {
		return nil, errors.New("configuration empty / not loaded")
	}
	return config, nil
}
