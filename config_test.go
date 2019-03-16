package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	Address string        `yaml:"address"`
	Count   int           `yaml:"count"`
	Passive bool          `yaml:"passive" default:"true"`
	Period  time.Duration `yaml:"period" default:"1m"`
}

var yml = `
---
address: http://example.com/
count: 23
passive: false
period: 2m22s
`
var badYml = "This is really NOT YAML."

// Load - default settings
func TestLoadDefault(t *testing.T) {
	cfg := Test{}
	Load("FUNAPP", "", &cfg)
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, true, cfg.Passive)
	assert.Equal(t, time.Minute, cfg.Period)
}

// FromYaml - good YAML
func TestFromYaml(t *testing.T) {
	cfg := Test{}
	err := FromYaml([]byte(yml), &cfg)
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://example.com/", cfg.Address)
	assert.Equal(t, 23, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, (2*time.Minute)+(22*time.Second), cfg.Period)
}

// FromYaml - bad YAML
func TestFromYamlBad(t *testing.T) {
	cfg := Test{}
	err := FromYaml([]byte(badYml), &cfg)
	assert.NotNil(t, err, "Expected an error because of bad YAML")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, cfg.Period, time.Duration(0))
}

// FromYamlFile - good YAML file
func TestFromYamlFile(t *testing.T) {
	file, err := ioutil.TempFile(".", "yaml_test")
	assert.Nil(t, err, "Got error trying to create temporary YAML file")
	assert.NotNil(t, file, "Failed to create temporary YAML file")
	defer os.Remove(file.Name())

	l, err := file.Write([]byte(yml))
	file.Close()
	assert.Nil(t, err, "Got error trying to write contents to temporary YAML file")
	assert.Equal(t, len(yml), l, "Mismatched bytes written to temporary YAML file")
	assert.FileExists(t, file.Name(), "Temporary YAML file doesn't exist")

	cfg := Test{}
	err = FromYamlFile(file.Name(), &cfg)
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://example.com/", cfg.Address)
	assert.Equal(t, 23, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, (2*time.Minute)+(22*time.Second), cfg.Period)
}

// FromYamlFile - bad YAML file
func TestFromYamlFileBad(t *testing.T) {
	file, err := ioutil.TempFile(".", "yaml_test_bad")
	assert.Nil(t, err, "Got error trying to create temporary YAML file")
	assert.NotNil(t, file, "Failed to create temporary YAML file")
	defer os.Remove(file.Name())

	l, err := file.Write([]byte(badYml))
	file.Close()
	assert.Nil(t, err, "Got error trying to write contents to temporary YAML file")
	assert.Equal(t, len(badYml), l, "Mismatched bytes written to temporary YAML file")
	assert.FileExists(t, file.Name(), "Temporary YAML file doesn't exist")

	cfg := Test{}
	err = FromYamlFile(file.Name(), &cfg)
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, cfg.Period, time.Duration(0))
}

// FromYamlFile - missing YAML file
func TestFromYamlFileNoFile(t *testing.T) {
	cfg := Test{}
	err := FromYamlFile("bogus_file_name", &cfg)
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, cfg.Period, time.Duration(0))
}

// FromEnvironment - env vars exist
func TestFromEnvironment(t *testing.T) {
	os.Setenv("FUNAPP_ADDRESS", "http://example.com/funapp")
	os.Setenv("FUNAPP_COUNT", "34")
	os.Setenv("FUNAPP_PASSIVE", "true")
	os.Setenv("FUNAPP_PERIOD", "11h11m11s")

	cfg := Test{}
	err := FromEnvironment("FUNAPP", &cfg)
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://example.com/funapp", cfg.Address)
	assert.Equal(t, 34, cfg.Count)
	assert.Equal(t, true, cfg.Passive)
	assert.Equal(t, (11*time.Hour)+(11*time.Minute)+(11*time.Second), cfg.Period)
}

// FromEnvironment - only required are set, rest should be default values
func TestFromEnvironmentDefault(t *testing.T) {
	os.Unsetenv("FUNAPP_ADDRESS")
	os.Unsetenv("FUNAPP_COUNT")
	os.Unsetenv("FUNAPP_PASSIVE")
	os.Unsetenv("FUNAPP_PERIOD")

	cfg := Test{}
	err := FromEnvironment("FUNAPP", &cfg)
	assert.Nil(t, err)
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, true, cfg.Passive)
	assert.Equal(t, cfg.Period, time.Minute)
}

type cfg struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	Sub     sub
}

type sub struct {
	Enabled bool `yaml:"enabled"`
	Level   int  `yaml:"level"`
}

func TestFromArguments(t *testing.T) {
	os.Args = []string{
		"Address=http://example.com",
		"Timeout=1h20m30s",
		"Sub.Enabled=true",
		"Sub.Level=42",
	}

	c := cfg{}
	err := FromArguments(os.Args, &c)
	assert.Nil(t, err)
	assert.Equal(t, c.Address, "http://example.com")
	assert.Equal(t, c.Timeout, (time.Hour)+(20*time.Minute)+(30*time.Second))
	assert.Equal(t, c.Sub.Enabled, true)
	assert.Equal(t, c.Sub.Level, 42)
}

// invalid parameter format is ignored
func TestFromArgumentsInvalidParameter(t *testing.T) {
	os.Args = []string{
		"Address:http://example.com",
	}

	c := cfg{}
	err := FromArguments(os.Args, &c)
	assert.Nil(t, err)
	assert.Equal(t, c.Address, "")
	assert.Equal(t, c.Timeout, time.Duration(0))
	assert.Equal(t, c.Sub.Enabled, false)
	assert.Equal(t, c.Sub.Level, 0)
}

// non-existant struct field name
func TestFromArgumentsInvalidFieldName(t *testing.T) {
	os.Args = []string{
		"Hey_dress=http://example.com",
	}

	c := cfg{}
	err := FromArguments(os.Args, &c)
	assert.NotNil(t, err)
	assert.Equal(t, c.Address, "")
	assert.Equal(t, c.Timeout, time.Duration(0))
	assert.Equal(t, c.Sub.Enabled, false)
	assert.Equal(t, c.Sub.Level, 0)
}

// non-parsable parameter value for the type
func TestFromArgumentsInvalidValue(t *testing.T) {
	os.Args = []string{
		"Timeout=SomeTime",
	}

	c := cfg{}
	err := FromArguments(os.Args, &c)
	assert.NotNil(t, err)
	assert.Equal(t, c.Address, "")
	assert.Equal(t, c.Timeout, time.Duration(0))
	assert.Equal(t, c.Sub.Enabled, false)
	assert.Equal(t, c.Sub.Level, 0)
}
