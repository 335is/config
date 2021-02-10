# Config

Loads configuration settings into a struct from the following sources in this order.

1. default specified in struct tag
1. config.yml YAML file
1. environment variables
1. command line parameters

Each source successively overrides the previous ones.

Nested structs are supported as each setting specifies its full path to the associated struct field.

## How to Use

See the examples folder and the unit tests for how to define and consume your configuration.

### Config Struct

The following config struct will contain our settings.
Note the default values defined in the struct tags.

```go
var cfg struct {
	Address string        `yaml:"address" default:"localhost"`
	Timeout time.Duration `yaml:"timeout" default:"5m"`
}
```

### YAML File

config.yml file overrides any struct defaults.

```yaml
---
address: http://example.com/
timeout: 1m
```

### Environment Variables

Environment variables override any matching config.yml file values.

```bash
export ADDRESS=http://example.com/funapp
export TIMEOUT=1m30s
```

### Command Line Arguments

Command line arguments override any matching environment variables.

```bash
./myapp address=http://example.com/home timeout=2m
```
