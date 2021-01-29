package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-lookup"
	yaml "gopkg.in/yaml.v2"
)

const (
	defaultConfigFile = "./config.yml"
)

// Load fills in the specified struct with configuration loaded from YAML, env vars, and command line arguments.
// It purposely ignores any errors from attempting to load from a specific source.
func Load(appName string, file string, struc interface{}) {
	if file == "" {
		file = defaultConfigFile
	}

	// overlay from local YAML config file
	FromYamlFile(file, struc)

	// overlay from environment variables
	FromEnvironment(appName, struc)

	// overlay from command line args
	FromArguments(os.Args[1:], struc)
}

// FromYaml extracts settings from a YAML string.
func FromYaml(yml []byte, struc interface{}) error {
	return yaml.Unmarshal(yml, struc)
}

// FromYamlFile extracts settings from a YAML file.
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
// We expect them to be named upper case, underscore separated, and prefixed with appName.
//	MYAPP_SERVER_ADDRESS=http://example.com
func FromEnvironment(appName string, struc interface{}) error {
	return envconfig.Process(strings.ToUpper(appName), struc)
}

// FromArguments extracts settings from a list of arguments such as those supplied on the command line.
// All args strings must be in the form key.path=value, where key.path is a period '.' or underscore '_' separated path to the struct member.
//	server.address=http://example.com
func FromArguments(args []string, struc interface{}) error {
	for _, arg := range args {
		kv := strings.Split(arg, "=")
		if len(kv) != 2 {
			continue
		}

		// transform period separators into underscores
		key := strings.ReplaceAll(kv[0], "_", ".")

		// find struct member matching key path
		value, err := lookup.LookupString(struc, key)
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
func ToYaml(struc interface{}) (string, error) {
	buff, err := yaml.Marshal(struc)
	if err == nil {
		return "---\n" + string(buff), nil
	}

	return "", err
}
