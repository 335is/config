# Config

Config loads configuration settings into a struct from the following sources in this order.

1. default struct values
1. config.yml YAML file
1. environment variables
1. command line parameters

Each source successively overrides the previous ones.

Nested structs are supported as each setting specifies its full path to the associated struct field.

## Requirements

### Go 1.12

[Download Go](https://golang.org/dl/)

### Dep 0.4.1

[Dep Releases](https://github.com/golang/dep/releases/)

## How to Use

See the examples folder and the unit tests for how to define and consume your configuration.

## Samples

### Default

The following cfg struct will contain our settings. Default values are coded in the tags.

```go
var cfg struct {
	Address string        `yaml:"address" default:"https://example.com"`
	Timeout time.Duration `yaml:"timeout" default:"30s"`
}
```

### YAML File

A config.yml file overrides defaults.

```yaml
---
address: http://example.com/
timeout: 1m
```

### Environment Variables

Environment variables override a config.yml file.

```bash
export MYAPP_ADDRESS=http://example.com/funapp
export MYAPP_TIMEOUT=1m30s
```

### Command Line Arguments

Command line arguments override environment variables.

```bash
./myapp address=http://example.com/home timeout=2m
```
