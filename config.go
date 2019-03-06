package config

import (
	"io/ioutil"
	"strings"

	"github.com/kelseyhightower/envconfig"
	yaml "gopkg.in/yaml.v2"
)

// Load fills in the specified struct with configuration loaded from YAML, env vars, etc.
// It purposely ignores any errors from attempting to load from a specific source.
func Load(appName string, struc interface{}) {
	// overlay from local YAML config file
	FromYamlFile("./config.yml", struc)

	// overlay from environment variables
	FromEnvironment(appName, struc)
}

// FromYaml extracts settings from a YAML string
func FromYaml(yml []byte, struc interface{}) error {
	return yaml.Unmarshal(yml, struc)
}

// FromYamlFile extracts settings from a YAML file
func FromYamlFile(path string, struc interface{}) error {
	// read YAML text file into a string
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// unmarshal from string to struct
	return FromYaml(yml, struc)
}

// FromEnvironment extracts settings from environment variables.
// We expect them to be named:
//		upper case
//		underscore separated
//		prefixed with appName
// For example:
//		MYAPP_CONSUL_ADDRESS
func FromEnvironment(appName string, struc interface{}) error {
	return envconfig.Process(strings.ToUpper(appName), struc)
}
