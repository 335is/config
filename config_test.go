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
	Passive bool          `yaml:"passive"`
	Period  time.Duration `yaml:"period"`
}

var yml = `---
address: http://example.com/
count: 23
passive: false
period: 2m22s
`
var notYml = "This is really NOT YAML."

// When there is no found YAML/ENV VARS/command line args, it should not override initialized struct members
func TestLoadEmptyConfigNoInitialization(t *testing.T) {
	cfg := Test{}
	Load("FUNAPP", "", &cfg)
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// When there is no found YAML/ENV VARS/command line args, it should not override initialized struct members
func TestLoadEmptyConfigWithInitialization(t *testing.T) {
	tm, _ := time.ParseDuration("1m11s")
	cfg := Test{
		Address: "http://google.com/",
		Count:   12,
		Passive: true,
		Period:  tm,
	}

	Load("FUNAPP", "", &cfg)
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://google.com/", cfg.Address)
	assert.Equal(t, 12, cfg.Count)
	assert.Equal(t, true, cfg.Passive)
	assert.Equal(t, tm, cfg.Period)
}

// good YAML string
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

// bad YAML string
func TestFromYamlBad(t *testing.T) {
	cfg := Test{}
	err := FromYaml([]byte(notYml), &cfg)
	assert.NotNil(t, err, "Expected an error because of bad YAML")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// good YAML file
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

// bad YAML file
func TestFromYamlFileBad(t *testing.T) {
	file, err := ioutil.TempFile(".", "yaml_test_bad")
	assert.Nil(t, err, "Got error trying to create temporary YAML file")
	assert.NotNil(t, file, "Failed to create temporary YAML file")
	defer os.Remove(file.Name())

	l, err := file.Write([]byte(notYml))
	file.Close()
	assert.Nil(t, err, "Got error trying to write contents to temporary YAML file")
	assert.Equal(t, len(notYml), l, "Mismatched bytes written to temporary YAML file")
	assert.FileExists(t, file.Name(), "Temporary YAML file doesn't exist")

	cfg := Test{}
	err = FromYamlFile(file.Name(), &cfg)
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// missing YAML file
func TestFromYamlFileNoFile(t *testing.T) {
	cfg := Test{}
	err := FromYamlFile("bogus_file_name", &cfg)
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// config struct to YAML string
func TestToYaml(t *testing.T) {
	tm, _ := time.ParseDuration("2m22s")
	cfg := Test{
		Address: "http://example.com/",
		Count:   23,
		Passive: false,
		Period:  tm,
	}

	s, err := ToYaml(&cfg)
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, yml, s)
}

// env vars exist
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

// only required are set, rest should be default values
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
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
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

// from command line arguments
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
	assert.Equal(t, "http://example.com", c.Address)
	assert.Equal(t, (time.Hour)+(20*time.Minute)+(30*time.Second), c.Timeout)
	assert.Equal(t, true, c.Sub.Enabled)
	assert.Equal(t, 42, c.Sub.Level)
}

// from command line arguments, that invalid parameter format is ignored
func TestFromArgumentsInvalidParameter(t *testing.T) {
	os.Args = []string{
		"Address:http://example.com",
	}

	c := cfg{}
	err := FromArguments(os.Args, &c)
	assert.Nil(t, err)
	assert.Equal(t, "", c.Address)
	assert.Equal(t, time.Duration(0), c.Timeout)
	assert.Equal(t, false, c.Sub.Enabled)
	assert.Equal(t, 0, c.Sub.Level)
}

// non-existant struct field name
func TestFromArgumentsInvalidFieldName(t *testing.T) {
	os.Args = []string{
		"Hey_dress=http://example.com",
	}

	c := cfg{}
	err := FromArguments(os.Args, &c)
	assert.NotNil(t, err)
	assert.Equal(t, "", c.Address)
	assert.Equal(t, time.Duration(0), c.Timeout)
	assert.Equal(t, false, c.Sub.Enabled)
	assert.Equal(t, 0, c.Sub.Level)
}

// non-parsable parameter value for the type
func TestFromArgumentsInvalidValue(t *testing.T) {
	os.Args = []string{
		"Timeout=SomeTime",
	}

	c := cfg{}
	err := FromArguments(os.Args, &c)
	assert.NotNil(t, err)
	assert.Equal(t, "", c.Address)
	assert.Equal(t, time.Duration(0), c.Timeout)
	assert.Equal(t, false, c.Sub.Enabled)
	assert.Equal(t, 0, c.Sub.Level)
}
