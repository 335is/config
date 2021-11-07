package config

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestInitialization struct {
	Address string
	Count   int
	Passive bool
	Period  time.Duration
}

type TestDefaults struct {
	Address string        `default:"http://stonewall.net"`
	Count   int           `default:"27"`
	Passive bool          `default:"true"`
	Period  time.Duration `default:"20m20s"`
}

type TestDefaultsNested struct {
	Address string        `yaml:"address" default:"https://dosan.com/"`
	Timeout time.Duration `yaml:"timeout" default:"11m11s"`
	Sub     SubNested
}

type SubNested struct {
	Enabled bool `yaml:"enabled" default:"true"`
	Level   int  `yaml:"level" default:"77"`
}

type TestYaml struct {
	Address string        `yaml:"address"`
	Count   int           `yaml:"count"`
	Passive bool          `yaml:"passive"`
	Period  time.Duration `yaml:"period"`
}

type TestEnvironment struct {
	Address string
	Count   int
	Passive bool
	Period  time.Duration
}

type TestArguments struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	Sub     sub
}

type sub struct {
	Enabled bool `yaml:"enabled"`
	Level   int  `yaml:"level"`
}

var yml = `---
address: http://example.com/
count: 23
passive: false
period: 2m22s
`
var notYml = "This is really NOT YAML."

// It should not overwrite uninitialized struct members if there is no matching YAML, environment variables, or command line arguments found.
func TestLoadEmptyConfigNoInitialization(t *testing.T) {
	cfg := TestInitialization{}
	Load("", &cfg)
	assert.NotNil(t, cfg)
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// It should not overwrite initialized struct members if there is no matching YAML, environment variables, or command line arguments found.
func TestLoadEmptyConfigWithInitialization(t *testing.T) {
	tm, _ := time.ParseDuration("1m11s")
	cfg := TestInitialization{
		Address: "http://google.com/",
		Count:   12,
		Passive: true,
		Period:  tm,
	}

	Load("", &cfg)
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://google.com/", cfg.Address)
	assert.Equal(t, 12, cfg.Count)
	assert.Equal(t, true, cfg.Passive)
	assert.Equal(t, tm, cfg.Period)
}

// It should set values of unitialized struct members from default struct tags.
func TestFromStructDefaultsNoInitialization(t *testing.T) {
	cfg := TestDefaults{}
	err := FromStructDefaults(&cfg)
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://stonewall.net", cfg.Address)
	assert.Equal(t, 27, cfg.Count)
	assert.Equal(t, true, cfg.Passive)
	assert.Equal(t, (20*time.Minute)+(20*time.Second), cfg.Period)
}

// It should not overwrite initialized struct members from default struct tags.
func TestFromStructDefaultsWithInitialization(t *testing.T) {
	tm, _ := time.ParseDuration("3m33s")
	cfg := TestDefaults{
		Address: "http://marble.net/",
		Count:   33,
		Passive: true,
		Period:  tm,
	}

	Load("", &cfg)
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://marble.net/", cfg.Address)
	assert.Equal(t, 33, cfg.Count)
	assert.Equal(t, true, cfg.Passive)
	assert.Equal(t, tm, cfg.Period)
}

func TestFromStructDefaultsNested(t *testing.T) {
	cfg := TestDefaultsNested{}
	err := FromStructDefaults(&cfg)

	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, "https://dosan.com/", cfg.Address)
	assert.Equal(t, (11*time.Minute)+(11*time.Second), cfg.Timeout)
	assert.Equal(t, true, cfg.Sub.Enabled)
	assert.Equal(t, 77, cfg.Sub.Level)
}

// parsing a good YAML string succeeds
func TestFromYaml(t *testing.T) {
	cfg := TestYaml{}
	err := FromYaml([]byte(yml), &cfg)
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://example.com/", cfg.Address)
	assert.Equal(t, 23, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, (2*time.Minute)+(22*time.Second), cfg.Period)
}

// parsing a bad YAML string fails
func TestFromYamlBad(t *testing.T) {
	cfg := TestYaml{}
	err := FromYaml([]byte(notYml), &cfg)
	assert.NotNil(t, err, "Expected an error because of bad YAML")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// parsing a good YAML file succeeds
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

	cfg := TestYaml{}
	err = FromYamlFile(file.Name(), &cfg)
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://example.com/", cfg.Address)
	assert.Equal(t, 23, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, (2*time.Minute)+(22*time.Second), cfg.Period)
}

// parsing a bad YAML file fails
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

	cfg := TestYaml{}
	err = FromYamlFile(file.Name(), &cfg)
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// trying to parse a missing YAML file fails
func TestFromYamlFileNoFile(t *testing.T) {
	cfg := TestYaml{}
	err := FromYamlFile("bogus_file_name", &cfg)
	assert.NotNil(t, err, "Expected an error")
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// parse a YAML string into a config struct
func TestToYaml(t *testing.T) {
	tm, _ := time.ParseDuration("2m22s")
	cfg := TestYaml{
		Address: "http://example.com/",
		Count:   23,
		Passive: false,
		Period:  tm,
	}

	s, err := ToYaml(&cfg)
	assert.Nil(t, err, "Expected no error")
	assert.Equal(t, yml, s)
}

// parse environment variables into config struct members
func TestFromEnvironment(t *testing.T) {
	os.Setenv("ADDRESS", "http://example.com/funapp")
	os.Setenv("COUNT", "34")
	os.Setenv("PASSIVE", "true")
	os.Setenv("PERIOD", "11h11m11s")

	cfg := TestEnvironment{}
	err := FromEnvironment(&cfg)
	assert.Nil(t, err, "Expected no error")
	assert.NotNil(t, cfg)
	assert.Equal(t, "http://example.com/funapp", cfg.Address)
	assert.Equal(t, 34, cfg.Count)
	assert.Equal(t, true, cfg.Passive)
	assert.Equal(t, (11*time.Hour)+(11*time.Minute)+(11*time.Second), cfg.Period)
}

// verify no matching environment variables leaves struct members untouched
func TestFromEnvironmentDefault(t *testing.T) {
	os.Unsetenv("ADDRESS")
	os.Unsetenv("COUNT")
	os.Unsetenv("PASSIVE")
	os.Unsetenv("PERIOD")

	cfg := TestEnvironment{}
	err := FromEnvironment(&cfg)
	assert.Nil(t, err)
	assert.Equal(t, "", cfg.Address)
	assert.Equal(t, 0, cfg.Count)
	assert.Equal(t, false, cfg.Passive)
	assert.Equal(t, time.Duration(0), cfg.Period)
}

// parse from command line arguments into struct members
func TestFromArguments(t *testing.T) {
	os.Args = []string{
		"Address=http://example.com",
		"Timeout=1h20m30s",
		"Sub.Enabled=true",
		"Sub.Level=42",
	}

	c := TestArguments{}
	err := FromArguments(os.Args, &c)
	assert.Nil(t, err)
	assert.Equal(t, "http://example.com", c.Address)
	assert.Equal(t, (time.Hour)+(20*time.Minute)+(30*time.Second), c.Timeout)
	assert.Equal(t, true, c.Sub.Enabled)
	assert.Equal(t, 42, c.Sub.Level)
}

// parse from command line arguments, verify that an invalid command line parameter (: instead of =) is ignored
func TestFromArgumentsInvalidParameter(t *testing.T) {
	os.Args = []string{
		"Address:http://example.com",
	}

	c := TestArguments{}
	err := FromArguments(os.Args, &c)
	assert.Nil(t, err)
	assert.Equal(t, "", c.Address)
	assert.Equal(t, time.Duration(0), c.Timeout)
	assert.Equal(t, false, c.Sub.Enabled)
	assert.Equal(t, 0, c.Sub.Level)
}

// verify that a non-matching command line parameter is ignored
func TestFromArgumentsInvalidFieldName(t *testing.T) {
	os.Args = []string{
		"Red_dress=http://example.com",
	}

	c := TestArguments{}
	err := FromArguments(os.Args, &c)
	assert.NotNil(t, err)
	assert.Equal(t, "", c.Address)
	assert.Equal(t, time.Duration(0), c.Timeout)
	assert.Equal(t, false, c.Sub.Enabled)
	assert.Equal(t, 0, c.Sub.Level)
}

// verify an invalid (non-parsable for the type) command line parameter value is ignored
func TestFromArgumentsInvalidValue(t *testing.T) {
	os.Args = []string{
		"Timeout=SomeTime",
	}

	c := TestArguments{}
	err := FromArguments(os.Args, &c)
	assert.NotNil(t, err)
	assert.Equal(t, "", c.Address)
	assert.Equal(t, time.Duration(0), c.Timeout)
	assert.Equal(t, false, c.Sub.Enabled)
	assert.Equal(t, 0, c.Sub.Level)
}
