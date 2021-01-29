# Config

Config loads configuration settings into a struct from the following sources in this order.

1. config.yml YAML file
1. environment variables
1. command line parameters

Each source successively overrides the previous ones.

Nested structs are supported as each setting specifies its full path to the associated struct field.

## Requirements

### Go 1.15

## How to Use

See the examples folder and the unit tests for how to define and consume your configuration.

## Samples

### Struct

The following cfg struct will contain our settings.

```go
var cfg struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}
```

### YAML File

A config.yml file overrides.

```yaml
---
address: http://example.com/
timeout: 1m
```

### Environment Variables

Environment variables overrides.

```bash
export MYAPP_ADDRESS=http://example.com/funapp
export MYAPP_TIMEOUT=1m30s
```

### Command Line Arguments

Command line arguments overrides.

```bash
./myapp address=http://example.com/home timeout=2m
```
