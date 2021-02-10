package config

import (
	"io/ioutil"
	"os"
	"strings"

	defaults "github.com/mcuadros/go-defaults"
	"github.com/mcuadros/go-lookup"
	"github.com/vrischmann/envconfig"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultConfigFile = "./config.yml"
)

// Load fills in the specified struct with configuration loaded from YAML, env vars, and command line arguments.
// It purposely ignores any errors from attempting to load from a specific source.
func Load(file string, v interface{}) {
	if file == "" {
		file = defaultConfigFile
	}

	// initialize with any "default:" struct tag values
	FromStructDefaults(v)

	// overlay from local YAML config file
	FromYamlFile(file, v)

	// overlay from environment variables
	FromEnvironment(v)

	// overlay from command line args
	FromArguments(os.Args[1:], v)
}

// FromStructDefaults initializes struct members from "default:" struc tags
func FromStructDefaults(v interface{}) error {
	defaults.SetDefaults(v)
	return nil
}

// FromYaml extracts settings from a YAML string.
func FromYaml(yml []byte, v interface{}) error {
	return yaml.Unmarshal(yml, v)
}

// FromYamlFile extracts settings from a YAML file.
func FromYamlFile(path string, v interface{}) error {
	// read YAML text file into a string
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// unmarshal from string to struct
	return FromYaml(yml, v)
}

// FromEnvironment extracts settings from environment variables.
// We expect them to be named upper case and underscore separated.
//	SERVER_ADDRESS=http://example.com
func FromEnvironment(v interface{}) error {
	return envconfig.InitWithOptions(v, envconfig.Options{AllOptional: true})
}

// FromArguments extracts settings from a list of arguments such as those supplied on the command line.
// All args strings must be in the form key.path=value, where key.path is a period '.' or underscore '_' separated path to the struct member.
//	server.address=http://example.com
func FromArguments(args []string, v interface{}) error {
	for _, arg := range args {
		kv := strings.Split(arg, "=")
		if len(kv) != 2 {
			continue
		}

		// transform period separators into underscores
		key := strings.ReplaceAll(kv[0], "_", ".")

		// find struct member matching key path
		value, err := lookup.LookupString(v, key)
		if err != nil {
			return err
		}

		// unmarshal the string into the struct field
		err = UnmarshalValue(kv[1], value)
		if err != nil {
			return err
		}
	}

	return nil
}

// ToYaml marshals the struc into a YAML string.
func ToYaml(v interface{}) (string, error) {
	buff, err := yaml.Marshal(v)
	if err == nil {
		return "---\n" + string(buff), nil
	}

	return "", err
}
